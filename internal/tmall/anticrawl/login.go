package anticrawl

import (
	"fmt"
	"strings"

	"github.com/tebeka/selenium"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/anticrawl"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

// GetTmallCookieStr
// @param usernameStr
// @param passwordStr
// @return cookieStr
// @return err
// @author: Kevineluo
func GetTmallCookieStr(usernameStr string, passwordStr string) (cookieStr string, err error) {
	wd, err := anticrawl.InitWebDriver()
	if err != nil {
		log.GLogger.Errorf("InitWebDriver failed, error: %s", err.Error())
	}
	defer wd.Quit()

	if err = wd.Get("https://login.tmall.com/"); err != nil {
		log.GLogger.Errorf("get page failed, error: %s", err.Error())
		return
	}
	// 登录框在嵌套的iframe中，需要切过去
	wd.SwitchFrame("J_loginIframe")
	// 等待element加载完全
	if err = wd.Wait(anticrawl.CheckDisplayed(selenium.ByCSSSelector, "#fm-login-id")); err != nil {
		log.GLogger.Errorf(err.Error())
		return
	}
	if err = wd.Wait(anticrawl.CheckDisplayed(selenium.ByCSSSelector, "#fm-login-password")); err != nil {
		log.GLogger.Errorf(err.Error())
		return
	}
	if err = wd.Wait(anticrawl.CheckDisplayed(selenium.ByCSSSelector, "#login-form > div.fm-btn > button")); err != nil {
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
