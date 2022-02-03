package anticrawl

import (
	"fmt"
	"strings"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func Displayed(by, elementName string) func(selenium.WebDriver) (bool, error) {
	return func(wd selenium.WebDriver) (ok bool, err error) {
		var el selenium.WebElement
		el, err = wd.FindElement(by, elementName)
		if err != nil {
			return
		}
		ok, err = el.IsDisplayed()
		return
	}
}

func GetCookieStr(usernameStr string, passwordStr string) (cookieStr string) {
	var (
		driverPath = "chromedriver"
		port       = 9222
	)

	service, err := selenium.NewChromeDriverService(driverPath, port, []selenium.ServiceOption{}...)
	if nil != err {
		fmt.Println("start a chromedriver service failed", err.Error())
		return
	}
	defer func() {
		_ = service.Stop()
	}()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--disable-gpu-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36",
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		fmt.Println("connect to the webDriver failed", err.Error())
		return
	}
	defer func() {
		_ = wd.Quit()
	}()
	err = wd.Get("https://passport.weibo.cn/signin/login?entry=mweibo&r=https://weibo.cn/")
	if err != nil {
		fmt.Println("get page failed", err.Error())
		return
	}
	err = wd.Wait(Displayed(selenium.ByCSSSelector, "#loginName"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = wd.Wait(Displayed(selenium.ByCSSSelector, "#loginPassword"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = wd.Wait(Displayed(selenium.ByCSSSelector, "#loginAction"))
	if err != nil {
		fmt.Println(err)
		return
	}
	username, err := wd.FindElement(selenium.ByCSSSelector, "#loginName")
	if err != nil {
		fmt.Println("get username failed", err.Error())
		return
	}
	password, err := wd.FindElement(selenium.ByCSSSelector, "#loginPassword")
	if err != nil {
		fmt.Println("get username failed", err.Error())
		return
	}
	submit, err := wd.FindElement(selenium.ByCSSSelector, "#loginAction")
	if err != nil {
		fmt.Println("get username failed", err.Error())
		return
	}
	err = username.SendKeys(usernameStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = password.SendKeys(passwordStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = submit.Click()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = wd.Wait(func(wdTemp selenium.WebDriver) (b bool, e error) {
		tit, err := wdTemp.Title()
		if err != nil {
			return false, nil
		}
		if tit != "我的首页" {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		fmt.Println(err)
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
