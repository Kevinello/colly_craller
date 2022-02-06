package colly

import (
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

var (
	Collector  *colly.Collector
	CrawlQueue *queue.Queue
)

// init
// @author: Kevineluo
func init() {

}

// InitCollector
// @return err
// @author: Kevineluo
func InitCollector() (err error) {
	// 初始化HomePageCollector
	Collector = colly.NewCollector(
		colly.CacheDir("./cache"),
		colly.Debugger(log.GLogger),
	)
	// 限制爬取速率
	Collector.Limit(&colly.LimitRule{
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})
	// 设置基础反爬插件
	extensions.RandomUserAgent(Collector)
	extensions.Referer(Collector)
	// 找到Category
	Collector.OnRequest(HandlerPrintRequestUrl)
	Collector.OnResponse(HandlerPrintResponseUrl)
	return
}

// InitCategoryQueue
// @return err
// @author: Kevineluo
func InitCategoryQueue() (err error) {
	// create a request queue with 2 consumer threads
	CrawlQueue, err = queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 100}, // Use default queue storage
	)
	return
}
