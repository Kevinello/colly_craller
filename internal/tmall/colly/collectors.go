package colly

import (
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

var (
	HomePageCollector  *colly.Collector
	CategoryCrawlQueue *queue.Queue
)

// init
// @author: Kevineluo
func init() {
	InitHomePageCollector()
	InitCategoryQueue()
}

// InitHomePageCollector
// @return err
// @author: Kevineluo
func InitHomePageCollector() (err error) {
	// 初始化HomePageCollector
	HomePageCollector = colly.NewCollector(
		colly.AllowedDomains("tmall.com", "www.tmall.com"),
		// colly.CacheDir("./cache"),
		colly.Debugger(log.GLogger),
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
	HomePageCollector.OnHTML(`script:contains(window\.\$data)`, HandlerGetCategoryUrl)
	HomePageCollector.OnRequest(HandlerSetCookie)
	HomePageCollector.OnResponse(HandlerPrintResponseUrl)
	return
}

// InitCategoryQueue
// @return err
// @author: Kevineluo
func InitCategoryQueue() (err error) {
	// create a request queue with 2 consumer threads
	CategoryCrawlQueue, err = queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 100}, // Use default queue storage
	)
	return
}
