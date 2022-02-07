package colly

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

var ()

// init
// @author: Kevineluo
func init() {

}

// InitCollector
// @return err
// @author: Kevineluo
func InitCollector() (Collector *colly.Collector) {
	// 初始化HomePageCollector
	Collector = colly.NewCollector(
		colly.CacheDir("./cache"),
		colly.Debugger(log.GLogger),
	)

	// 设置基础反爬插件
	extensions.RandomUserAgent(Collector)
	extensions.Referer(Collector)
	// 默认Handler
	Collector.OnRequest(HandlerPrintRequestUrl)
	Collector.OnResponse(HandlerPrintResponseUrl)
	return
}

// InitCrawlQueue
// @return err
// @author: Kevineluo
func InitCrawlQueue(threadNum int, storage queue.Storage) (CrawlQueue *queue.Queue, err error) {
	// create a request queue with 2 consumer threads
	CrawlQueue, err = queue.New(
		2,       // Number of consumer threads
		storage, // storage
	)
	return
}
