package anticrawl

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

var (
	CHROMEDRIVER_PATH        = pkg.GetEnv("CHROMEDRIVER_PATH", "/data/service/chromedriver")
	CHROME_BINARY_PATH       = pkg.GetEnv("CHROME_BINARY_PATH", "/root/colly_crawler/vendor/chrome-linux/chrome")
	SELENIUM_SERVICE_PORT, _ = strconv.Atoi(pkg.GetEnv("SELENIUM_SERVICE_PORT", "9000"))

	SeleniumService *selenium.Service
)

func init() {
	var err error
	opts := []selenium.ServiceOption{
		selenium.Output(os.Stderr),
	}
	selenium.SetDebug(true)

	SeleniumService, err = selenium.NewChromeDriverService(CHROMEDRIVER_PATH, SELENIUM_SERVICE_PORT, opts...)
	if err != nil {
		fmt.Println("start chromedriver service failed", err.Error())
		return
	}
}

// Displayed
// @param by
// @param elementName
// @return selenium.WebDriver
// @return func(selenium.WebDriver) (bool, error)
// @author: Kevineluo
func Displayed(by, elementName string) func(selenium.WebDriver) (bool, error) {
	return func(wd selenium.WebDriver) (ok bool, err error) {
		el, _ := wd.FindElement(by, elementName)
		if el != nil {
			ok, _ = el.IsDisplayed()
		}
		return
	}
}

// GetCookieStr
// @param usernameStr
// @param passwordStr
// @return cookieStr
// @return err
// @author: Kevineluo
func GetCookieStr(usernameStr string, passwordStr string) (cookieStr string, err error) {
	caps := selenium.Capabilities{"browserName": "chrome"}
	// prefs 禁止图片加载，加快渲染速度
	prefs := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	// chromeCaps 设置浏览器参数，随机生成user-agent
	chromeCaps := chrome.Capabilities{
		Prefs: prefs,
		Path:  CHROME_BINARY_PATH,
		// 设置为开发者模式，防止被各大网站识别出来使用了Selenium
		ExcludeSwitches: []string{"enable-automation"},
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--disable-gpu-sandbox",
			fmt.Sprintf("--user-agent=%s", RandomUserAgent()),
		},
	}
	caps.AddChrome(chromeCaps)
	// Connect to the WebDriver instance running locally.
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", SELENIUM_SERVICE_PORT))
	if err != nil {
		log.GLogger.Errorf("connect to the webDriver failed, error: %s", err.Error())
		return
	}
	defer wd.Quit()

	if err = wd.Get("https://login.tmall.com/"); err != nil {
		log.GLogger.Errorf("get page failed, error: %s", err.Error())
		return
	}
	// 登录框在嵌套的iframe中，需要切过去
	wd.SwitchFrame("J_loginIframe")
	// 等待element加载完全
	if err = wd.Wait(Displayed(selenium.ByCSSSelector, "#fm-login-id")); err != nil {
		log.GLogger.Errorf(err.Error())
		return
	}
	if err = wd.Wait(Displayed(selenium.ByCSSSelector, "#fm-login-password")); err != nil {
		log.GLogger.Errorf(err.Error())
		return
	}
	if err = wd.Wait(Displayed(selenium.ByCSSSelector, "#login-form > div.fm-btn > button")); err != nil {
		log.GLogger.Errorf(err.Error())
		return
	}
	// 获取对应element
	username, err := wd.FindElement(selenium.ByCSSSelector, "#fm-login-id")
	if err != nil {
		log.GLogger.Errorf("get username failed, error: %s", err.Error())
		return
	}
	password, err := wd.FindElement(selenium.ByCSSSelector, "#fm-login-password")
	if err != nil {
		log.GLogger.Errorf("get password failed, error: %s", err.Error())
		return
	}
	login, err := wd.FindElement(selenium.ByCSSSelector, "#login-form > div.fm-btn > button")
	if err != nil {
		log.GLogger.Errorf("get login button failed, error: %s", err.Error())
		return
	}
	// 填入用户名密码后点击登录
	err = username.SendKeys(usernameStr)
	if err != nil {
		log.GLogger.Error(err.Error())
		return
	}
	err = password.SendKeys(passwordStr)
	if err != nil {
		log.GLogger.Error(err.Error())
		return
	}
	err = login.Click()
	if err != nil {
		log.GLogger.Error(err.Error())
		return
	}
	// 确认是否跳转到首页
	err = wd.Wait(func(wdTemp selenium.WebDriver) (isTargetPage bool, err error) {
		title, err := wdTemp.Title()
		if err != nil {
			return
		}
		if !strings.Contains(title, "理想生活上天猫") {
			err = fmt.Errorf("not the target page!!!")
			return
		}
		isTargetPage = true
		return
	})
	if err != nil {
		log.GLogger.Error(err.Error())
		return
	}

	cookieLst, err := wd.GetCookies()
	if err != nil {
		fmt.Println(err)
		return
	}
	var cookieArr []string
	for _, c := range cookieLst {
		cookieArr = append(cookieArr, c.Name+"="+c.Value)
	}
	cookieStr = strings.Join(cookieArr, ";")
	log.GLogger.Logger.Debugf("Get cookieStr: %s", cookieStr)
	return
}
