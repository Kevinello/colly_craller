package anticrawl

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

var (
	CHROMEDRIVER_PATH        = pkg.GetEnv("CHROMEDRIVER_PATH", "/data/service/chromedriver")
	CHROME_BINARY_PATH       = pkg.GetEnv("CHROME_BINARY_PATH", "/root/colly_crawler/vendor/chrome-linux/chrome")
	SELENIUM_SERVICE_PORT, _ = strconv.Atoi(pkg.GetEnv("SELENIUM_SERVICE_PORT", "9000"))
)

func init() {

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
	opts := []selenium.ServiceOption{
		selenium.Output(os.Stderr),
	}
	selenium.SetDebug(true)

	service, err := selenium.NewChromeDriverService(CHROMEDRIVER_PATH, SELENIUM_SERVICE_PORT, opts...)
	if err != nil {
		fmt.Println("start chromedriver service failed", err.Error())
		return
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	// imagCaps 禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	// chromeCaps 设置浏览器参数，随机生成user-agent
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  CHROME_BINARY_PATH,
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
	// 等待element加载完全
	if err = wd.WaitWithTimeoutAndInterval(Displayed(selenium.ByCSSSelector, "#fm-login-id"), time.Second*5, time.Second); err != nil {
		log.GLogger.Errorf(err.Error())
		return
	}
	if err = wd.WaitWithTimeoutAndInterval(Displayed(selenium.ByCSSSelector, "#fm-login-password"), time.Second*5, time.Second); err != nil {
		log.GLogger.Errorf(err.Error())
		return
	}
	if err = wd.WaitWithTimeoutAndInterval(Displayed(selenium.ByCSSSelector, "#login-form > div.fm-btn > button"), time.Second*5, time.Second); err != nil {
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
	login, err := wd.FindElement(selenium.ByCSSSelector, "button.fm-btn:contain(登录)")
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
		if title != "天猫tmall.com--理想生活上天猫" {
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
	// 确认是否登陆成功
	if err = wd.Wait(Displayed(selenium.ByCSSSelector, "#login-info1")); err != nil {
		log.GLogger.Errorf(err.Error())
		return
	}
	if err = wd.Wait(Displayed(selenium.ByCSSSelector, "span:contain(Hi)")); err != nil {
		log.GLogger.Errorf(err.Error())
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
	cookieStr = strings.Join(cookieArr, "; ")
	return
}
