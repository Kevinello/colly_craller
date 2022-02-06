package anticrawl

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg"
)

var (
	CHROMEDRIVER_PATH        = pkg.GetEnv("CHROMEDRIVER_PATH", "/data/service/chromedriver")
	CHROME_BINARY_PATH       = pkg.GetEnv("CHROME_BINARY_PATH", "/root/colly_crawler/vendor/chrome-linux/chrome")
	SELENIUM_SERVICE_PORT, _ = strconv.Atoi(pkg.GetEnv("SELENIUM_SERVICE_PORT", "9000"))

	SeleniumService *selenium.Service
)

// init
// @author: Kevineluo
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

// InitWebDriver
// 注意defer wd.Quit()关闭webdriver
// @return wd
// @return err
// @author: Kevineluo
func InitWebDriver() (wd selenium.WebDriver, err error) {
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
	wd, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", SELENIUM_SERVICE_PORT))
	if err != nil {
		return
	}
	return
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
