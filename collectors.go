package collycrawller

import (
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var (
	HomePageCollector *colly.Collector
)

func init() {
	// 初始化Collecctor
	HomePageCollector = colly.NewCollector(
		colly.AllowedDomains("tmall.com", "www.tmall.com"),
		colly.CacheDir("./cache"),
		colly.Debugger(GLogger),
	)
	// 限制爬取速率
	HomePageCollector.Limit(&colly.LimitRule{
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})
	// 设置基础反爬插件
	extensions.RandomUserAgent(HomePageCollector)
	extensions.Referer(HomePageCollector)
	// 找到Category
	HomePageCollector.OnHTML("div[class~=Category--category]", func(h *colly.HTMLElement) {
		GLogger.Info("Category found")
		GLogger.Info(h.Text)
	})
	HomePageCollector.OnRequest(HandlerPrintUrl)
}
