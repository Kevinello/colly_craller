package web

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tebeka/selenium"
	"kevinello.ltd/kevinello/collycrawler/internal/jd/colly"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/anticrawl"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

const (
	searchUrlFmt = "https://search.jd.com/Search?keyword=%s&page=%d"
)

var (
	// ItemUrlChan 非缓冲Url Channel，通过阻塞限制爬取速率
	ItemUrlChan = make(chan string)

	// ITEM_CRAWL_INTERVAL item爬取间隔
	ITEM_CRAWL_INTERVAL, _ = strconv.Atoi(pkg.GetEnv("ITEM_CRAWL_INTERVAL", "1000"))
)

func init() {
	go ItemUrlHandler(ItemUrlChan)
}

func ItemUrlHandler(itemUrlChan chan string) {
	for {
		select {
		case url := <-itemUrlChan:
			log.GLogger.Infof("Get item url: %s", url)
			colly.ItemQueue.AddURL(url)
			// 限制爬取速率
			time.Sleep(time.Duration(ITEM_CRAWL_INTERVAL) * time.Millisecond)
		}
	}
}

func CollectItemUrl(keyword string, urlChan chan string, pageNum int) {
	// 最后关闭SeleniumService
	defer anticrawl.SeleniumService.Stop()

	wd, err := anticrawl.InitWebDriver()
	if err != nil {
		log.GLogger.Errorf("InitWebDriver failed, error: %s", err.Error())
		return
	}
	defer wd.Quit()

	for i := 0; i < pageNum; i++ {
		searchUrl := fmt.Sprintf(searchUrlFmt, keyword, 2*i+1)
		if err = wd.Get(searchUrl); err != nil {
			log.GLogger.Errorf("get page failed, error: %s", err.Error())
			return
		}

		// 等待element加载完全
		if err = wd.Wait(anticrawl.CheckDisplayed(selenium.ByCSSSelector, "#J_goodsList")); err != nil {
			log.GLogger.Errorf(err.Error())
			return
		}

		body, err := wd.FindElement(selenium.ByCSSSelector, "body")
		if err != nil {
			log.GLogger.Errorf("get body failed, error: %s", err.Error())
			return
		}

		// 发送空格使页面滑到底
		for i := 0; i < 30; i++ {
			body.SendKeys(selenium.EndKey)
		}
		itemList, err := wd.FindElements(selenium.ByCSSSelector, "div.gl-i-wrap > div.p-img > a")
		if err != nil {
			log.GLogger.Errorf("get itemList failed, error: %s", err.Error())
			return
		}
		for _, item := range itemList {
			itemUrl, err := item.GetAttribute("href")
			if err != nil {
				log.GLogger.Errorf("get itemUrl failed, error: %s", err.Error())
				continue
			}
			urlChan <- itemUrl
		}
	}
}

func CheckLoadMore(by, elementName string) func(selenium.WebDriver) (bool, error) {
	return func(wd selenium.WebDriver) (ok bool, err error) {
		el, _ := wd.FindElement(by, elementName)
		if el != nil {
			ok, _ = el.IsDisplayed()
		}
		if ok {
			var (
				body selenium.WebElement
				pos  *selenium.Point
			)
			pos, err = el.Location()
			if err != nil {
				log.GLogger.Errorf("get body failed, error: %s", err.Error())
				return
			}
			body, err = wd.FindElement(selenium.ByCSSSelector, "body")
			if err != nil {
				log.GLogger.Errorf("get body failed, error: %s", err.Error())
				return
			}
			body.MoveTo(pos.X, pos.Y)
		}
		return
	}
}
