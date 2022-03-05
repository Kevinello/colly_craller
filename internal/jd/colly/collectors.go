package colly

import (
	"time"

	"github.com/gocolly/colly"
	_colly "kevinello.ltd/kevinello/collycrawler/internal/pkg/colly"
)

var (
	ItemCollector          *colly.Collector
	WareBussinessCollector *colly.Collector
	CommentCollector       *colly.Collector
)

func init() {
	InitItemCollector()
	InitWareBussinessCollector()
	InitCommentCollector()
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
