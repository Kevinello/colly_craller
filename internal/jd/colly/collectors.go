package colly

import (
	"os"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	_colly "kevinello.ltd/kevinello/collycrawler/internal/pkg/colly"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

var (
	ItemCollector *colly.Collector
	ItemQueue     *queue.Queue
)

func init() {
	InitItemCollector()
	if err := InitItemQueue(); err != nil {
		log.GLogger.Errorf("InitItemQueue error: %s", err.Error())
		os.Exit(1)
	}
}

func InitItemCollector() {
	ItemCollector = _colly.InitCollector()
	// 限制爬取速率
	ItemCollector.Limit(&colly.LimitRule{
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})
	// 设置url过滤
	ItemCollector.AllowedDomains = []string{
		"item.jd.com",
	}
	ItemCollector.OnRequest(HandlerFindItemIdFromUrl)
	ItemCollector.OnHTML("div.itemInfo-wrap", HandlerParseItemInfo)
}

func InitItemQueue() (err error) {
	ItemQueue, err = _colly.InitCrawlQueue(2, &queue.InMemoryQueueStorage{MaxSize: 1000})
	return
}
