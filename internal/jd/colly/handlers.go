package colly

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"kevinello.ltd/kevinello/collycrawler/internal/jd/storage"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

const (
	priceUrlFormatter      = "https://p.3.cn/prices/mgets?skuIds=J_%s"
	commentUrlFormatter    = "https://club.jd.com/comment/productCommentSummaries.action?referenceIds=%s"
	wareBussinessFormatter = "https://item-soa.jd.com/getWareBusiness?&skuId=%s"
)

// HandlerFindItemIdFromUrl 从item url中获取item id
// @param r
// @author: Kevineluo
func HandlerFindItemIdFromUrl(r *colly.Request) {
	log.GLogger.Infof("visiting %s", r.URL.String())

	// 从url获取itemID
	itemID, err := GetItemIDFromUrl(r.URL.String())
	if err != nil {
		log.GLogger.Errorf(err.Error())
		return
	}
	log.GLogger.Infof("find item_id: %s", itemID)

	// 存储ItemID到ItemStorageMap
	saveRes := make(chan int)
	itemSaveMessage := &storage.ItemSaveMessage{
		ItemID:  itemID,
		SaveRes: saveRes,
	}
	storage.ItemStorageChan <- itemSaveMessage
	// 阻塞获取存储状态
	status := <-saveRes
	log.GLogger.Infof("item save status: %d", status)

	// 在collector的上下文中放入item_id
	PriceCollector.OnRequest(func(r *colly.Request) { r.Ctx.Put("item_id", itemID) })
	priceUrl := fmt.Sprintf(priceUrlFormatter, itemID)
	PriceCollector.Visit(priceUrl)

	// 在collector的上下文中放入item_id
	CommentCollector.OnRequest(func(r *colly.Request) { r.Ctx.Put("item_id", itemID) })
	commentUrl := fmt.Sprintf(commentUrlFormatter, itemID)
	CommentCollector.Visit(commentUrl)
}

func HandlerCollectSkuNum(h *colly.HTMLElement) {
	skuSelection := h.DOM.Children()
	log.GLogger.Infof("sku num: %d", skuSelection.Length())
	// 存储ItemID到ItemStorageMap
	itemID := h.Request.Ctx.Get("item_id")
	saveRes := make(chan int)
	itemSaveMessage := &storage.ItemSaveMessage{
		ItemID:    itemID,
		SaveField: "SkuNum",
		SaveValue: skuSelection.Length(),
		SaveRes:   saveRes,
	}
	storage.ItemStorageChan <- itemSaveMessage
	// 阻塞获取存储状态
	status := <-saveRes
	log.GLogger.Infof("item save status: %d", status)
}

// HandlerCollectPrice 从price接口收集商品价格
// @param r
// @author: Kevineluo
func HandlerCollectPrice(r *colly.Response) {
	// log.GLogger.Debugf("Get response: %s", string(r.Body))
	jsonStr := string(r.Body)
	jsonStr = strings.TrimSpace(jsonStr)

	priceResponses := storage.PriceResponse{}
	err := json.Unmarshal([]byte(jsonStr), &priceResponses)
	if err != nil {
		log.GLogger.Errorf("error when Unmarshal PriceResponse of Request[%s]: %s", r.Request.URL, err.Error())
		return
	}
	if len(priceResponses) > 0 {
		log.GLogger.Infof("Collect price from Request[%s]: %+v", r.Request.URL, priceResponses[0])
	} else {
		log.GLogger.Errorf("Can't find price from Request[%s]", r.Request.URL)
		return
	}
	priceResponse := priceResponses[0]
	price, err := strconv.Atoi(priceResponse.Price)
	if err != nil {
		log.GLogger.Errorf("illegal price[%s] from Request[%s]", priceResponse.Price, r.Request.URL)
		return
	}

	// 存储ItemID到ItemStorageMap
	itemID := r.Request.Ctx.Get("item_id")
	saveRes := make(chan int)
	itemSaveMessage := &storage.ItemSaveMessage{
		ItemID:    itemID,
		SaveField: "Price",
		SaveValue: price,
		SaveRes:   saveRes,
	}
	storage.ItemStorageChan <- itemSaveMessage
	// 阻塞获取存储状态
	status := <-saveRes
	log.GLogger.Infof("item save status: %d", status)
}

// HandlerCollectComment 从Comment接口收集评价信息
// @param r
// @author: Kevineluo
func HandlerCollectComment(r *colly.Response) {
	// log.GLogger.Debugf("Get response: %s", string(r.Body))
	jsonStr := string(r.Body)
	jsonStr = strings.TrimSpace(jsonStr)

	commentResponse := storage.CommentResponse{}
	err := json.Unmarshal([]byte(jsonStr), &commentResponse)
	if err != nil {
		log.GLogger.Errorf("error when Unmarshal CommentResponse of Request[%s]: %s", r.Request.URL, err.Error())
		return
	}
	if len(commentResponse.CommentsCount) > 0 {
		log.GLogger.Infof("Collect comment from Request[%s]: %+v", r.Request.URL, commentResponse.CommentsCount[0])
	} else {
		log.GLogger.Errorf("Can't find comment from Request[%s]", r.Request.URL)
		return
	}
	commentCount := commentResponse.CommentsCount[0]
	// TODO: 价格转换
	commentSummary := &storage.CommentSummary{
		AverageScore: commentCount.AverageScore,
		CommentCount: commentCount.CommentCount,
		GoodCount:    commentCount.GoodCount,
		GoodRate:     commentCount.GoodRate,
		GeneralCount: commentCount.GeneralCount,
		GeneralRate:  commentCount.GeneralRate,
		PoorCount:    commentCount.PoorCount,
		PoorRate:     commentCount.PoorRate,
	}

	// 存储ItemID到ItemStorageMap
	itemID := r.Request.Ctx.Get("item_id")
	saveRes := make(chan int)
	itemSaveMessage := &storage.ItemSaveMessage{
		ItemID:    itemID,
		SaveField: "CommentSummary",
		SaveValue: commentSummary,
		SaveRes:   saveRes,
	}
	storage.ItemStorageChan <- itemSaveMessage
	// 阻塞获取存储状态
	status := <-saveRes
	log.GLogger.Infof("item save status: %d", status)
}
