package web

import (
	"fmt"

	"github.com/tebeka/selenium"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/anticrawl"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

const (
	searchUrlFmt = "https://search.jd.com/Search?keyword=%s"
)

func CollectItemUrl(keyword string, urlChan chan string) {
	wd, err := anticrawl.InitWebDriver()
	if err != nil {
		log.GLogger.Errorf("InitWebDriver failed, error: %s", err.Error())
		return
	}
	defer wd.Quit()

	searchUrl := fmt.Sprintf(searchUrlFmt, keyword)
	if err = wd.Get(searchUrl); err != nil {
		log.GLogger.Errorf("get page failed, error: %s", err.Error())
		return
	}

	// 等待element加载完全
	if err = wd.Wait(anticrawl.Displayed(selenium.ByCSSSelector, "#J_goodsList")); err != nil {
		log.GLogger.Errorf(err.Error())
		return
	}

}
