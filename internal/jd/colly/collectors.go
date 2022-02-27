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
	ItemCollector          *colly.Collector
	WareBussinessCollector *colly.Collector
	CommentCollector       *colly.Collector
	ItemQueue              *queue.Queue
)

func init() {
	InitItemCollector()
	InitWareBussinessCollector()
	InitCommentCollector()

	// 初始化Item爬取队列
	if err := InitItemQueue(); err != nil {
		log.GLogger.Errorf("InitItemQueue error: %s", err.Error())
		os.Exit(1)
	}
}

// InitItemCollector 初始化ItemCollector
// @author: Kevineluo
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
	ItemCollector.OnHTML(`#choose-attr-1 > div.dd`, HandlerCollectSkuNum)
}

// InitWareBussinessCollector 初始化WareBussinessCollector
// @author: Kevineluo
func InitWareBussinessCollector() {
	WareBussinessCollector = _colly.InitCollector()
	// 限制爬取速率
	WareBussinessCollector.Limit(&colly.LimitRule{
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})
	// 设置url过滤
	WareBussinessCollector.AllowedDomains = []string{
		"item-soa.jd.com",
	}
	WareBussinessCollector.OnResponse(HandlerCollectWareBussiness)
}

// InitCommentCollector 初始化InitCommentCollector
// @author: Kevineluo
func InitCommentCollector() {
	CommentCollector = _colly.InitCollector()
	// 限制爬取速率
	CommentCollector.Limit(&colly.LimitRule{
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})
	// 设置url过滤
	CommentCollector.AllowedDomains = []string{
		"club.jd.com",
	}
	CommentCollector.OnResponse(HandlerCollectComment)
}

func InitItemQueue() (err error) {
	ItemQueue, err = _colly.InitCrawlQueue(2, &queue.InMemoryQueueStorage{MaxSize: 1000})
	return
}
