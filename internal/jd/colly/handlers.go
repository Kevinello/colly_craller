package colly

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/tebeka/selenium"
	"kevinello.ltd/kevinello/collycrawler/internal/jd/storage"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/anticrawl"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

const (
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
	// 在collector的上下文中放入item_id
	r.Ctx.Put("item_id", itemID)

	// 存储ItemID到ItemStorageMap
	saveRes := make(chan int)
	itemSaveMessage := &storage.ItemSaveMessage{
		ItemID:  itemID,
		SaveRes: saveRes,
	}
	storage.ItemStorageChan <- itemSaveMessage
	// 阻塞获取存储状态
	status := <-saveRes
	log.GLogger.Infof("item[%s] --- item save status: %d", itemID, status)

	// 限制爬取速率
	time.Sleep(time.Second)

	// 在collector的上下文中放入item_id
	WareBussinessCollector.OnRequest(func(r *colly.Request) { r.Ctx.Put("item_id", itemID) })
	wareBussinessUrl := fmt.Sprintf(wareBussinessFormatter, itemID)
	WareBussinessCollector.Visit(wareBussinessUrl)

	// 限制爬取速率
	time.Sleep(time.Second)

	// 在collector的上下文中放入item_id
	CommentCollector.OnRequest(func(r *colly.Request) { r.Ctx.Put("item_id", itemID) })
	commentUrl := fmt.Sprintf(commentUrlFormatter, itemID)
	CommentCollector.Visit(commentUrl)
}

// HandlerCollectSkuNum
// @param h
// @author: Kevineluo
func HandlerCollectSkuNum(h *colly.HTMLElement) {
	itemID := h.Request.Ctx.Get("item_id")
	skuSelection := h.DOM.Children()
	log.GLogger.Infof("item[%s] --- collect sku num: %d", itemID, skuSelection.Length())

	// 存储SkuNum到ItemStorageMap
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
	log.GLogger.Infof("item[%s] --- item save status: %d", itemID, status)
}

// HandlerCollectWareBussiness 从price接口收集商品价格
// @param r
// @author: Kevineluo
func HandlerCollectWareBussiness(r *colly.Response) {
	itemID := r.Request.Ctx.Get("item_id")
	jsonStr := string(r.Body)
	jsonStr = strings.TrimSpace(jsonStr)

	wareBussinessResponse := storage.WareBussinessResponse{}
	err := json.Unmarshal([]byte(jsonStr), &wareBussinessResponse)
	if err != nil {
		log.GLogger.Errorf("item[%s] --- error when Unmarshal wareBussinessResponse of Request[%s]: %s", itemID, r.Request.URL, err.Error())
		wd, err := anticrawl.InitWebDriver()
		if err != nil {
			log.GLogger.Errorf("item[%s] --- error when using selenium to InitWebDriver of Request[%s]: %s", itemID, r.Request.URL, err.Error())
			return
		}
		defer wd.Quit()
		// 进入页面
		if err = wd.Get(r.Request.URL.String()); err != nil {
			log.GLogger.Errorf("get page failed, error: %s", err.Error())
			return
		}
		// 等待element加载完全
		if err = wd.Wait(anticrawl.CheckDisplayed(selenium.ByCSSSelector, "body > pre")); err != nil {
			log.GLogger.Errorf(err.Error())
			return
		}
		response, err := wd.FindElement(selenium.ByCSSSelector, "body > pre")
		if err != nil {
			log.GLogger.Errorf("item[%s] --- error when using selenium to get response of Request[%s]: %s", itemID, r.Request.URL, err.Error())
			return
		}
		jsonStr, err = response.Text()
		if err != nil {
			log.GLogger.Errorf("item[%s] --- error when using selenium to get text of Request[%s]: %s", itemID, r.Request.URL, err.Error())
			return
		}
		jsonStr = strings.TrimSpace(jsonStr)

		err = json.Unmarshal([]byte(jsonStr), &wareBussinessResponse)
		if err != nil {
			log.GLogger.Errorf("item[%s] --- error when Unmarshal wareBussinessResponse of Request[%s]: %s", itemID, r.Request.URL, err.Error())
			return
		}
	}
	if !reflect.DeepEqual(wareBussinessResponse, storage.WareBussinessResponse{}) {
		log.GLogger.Debugf("item[%s] --- Collect wareBussiness from Request[%s]: %+v", itemID, r.Request.URL, wareBussinessResponse)
	} else {
		log.GLogger.Errorf("item[%s] --- Can't find wareBussiness from the response of Request[%s]", itemID, r.Request.URL)
		return
	}

	// 填充itemID
	wareBussinessResponse.Price.ItemID = itemID
	for idx := range wareBussinessResponse.Promotion.Activities {
		wareBussinessResponse.Promotion.Activities[idx].ItemID = itemID
	}
	wareBussinessResponse.ShopInfo.Shop.ItemID = itemID
	wareBussinessResponse.ShopInfo.CustomerService.ItemID = itemID

	saveRes := make(chan int)
	var status int
	// 存储Price到ItemStorageMap
	itemSaveMessage := &storage.ItemSaveMessage{
		ItemID:    itemID,
		SaveField: "Price",
		SaveValue: &wareBussinessResponse.Price,
		SaveRes:   saveRes,
	}
	storage.ItemStorageChan <- itemSaveMessage
	// 阻塞获取存储状态
	status = <-saveRes
	log.GLogger.Infof("item[%s] --- item save status: %d", itemID, status)
	// 存储Shop到ItemStorageMap
	itemSaveMessage = &storage.ItemSaveMessage{
		ItemID:    itemID,
		SaveField: "Shop",
		SaveValue: &wareBussinessResponse.ShopInfo.Shop,
		SaveRes:   saveRes,
	}
	storage.ItemStorageChan <- itemSaveMessage
	// 阻塞获取存储状态
	status = <-saveRes
	log.GLogger.Infof("item[%s] --- item save status: %d", itemID, status)
	// 存储CustomerService到ItemStorageMap
	itemSaveMessage = &storage.ItemSaveMessage{
		ItemID:    itemID,
		SaveField: "CustomerService",
		SaveValue: &wareBussinessResponse.ShopInfo.CustomerService,
		SaveRes:   saveRes,
	}
	storage.ItemStorageChan <- itemSaveMessage
	// 阻塞获取存储状态
	status = <-saveRes
	log.GLogger.Infof("item[%s] --- item save status: %d", itemID, status)
	// 存储Activities到ItemStorageMap
	itemSaveMessage = &storage.ItemSaveMessage{
		ItemID:    itemID,
		SaveField: "Activities",
		SaveValue: &wareBussinessResponse.Promotion.Activities,
		SaveRes:   saveRes,
	}
	storage.ItemStorageChan <- itemSaveMessage
	// 阻塞获取存储状态
	status = <-saveRes
	log.GLogger.Infof("item[%s] --- item save status: %d", itemID, status)
	log.GLogger.Infof("---------- wait for 10s ----------")
	// time.Sleep(10 * time.Second)
}

// HandlerCollectComment 从Comment接口收集评价信息
// @param r
// @author: Kevineluo
func HandlerCollectComment(r *colly.Response) {
	itemID := r.Request.Ctx.Get("item_id")
	jsonStr := string(r.Body)
	jsonStr = strings.TrimSpace(jsonStr)

	commentResponse := storage.CommentResponse{}
	err := json.Unmarshal([]byte(jsonStr), &commentResponse)
	if err != nil {
		log.GLogger.Errorf("item[%s] --- error when Unmarshal CommentResponse of Request[%s]: %s", itemID, r.Request.URL, err.Error())
		return
	}
	if len(commentResponse.CommentsCount) > 0 {
		log.GLogger.Debugf("item[%s] --- Collect comment from Request[%s]: %+v", itemID, r.Request.URL, commentResponse.CommentsCount[0])
	} else {
		log.GLogger.Errorf("item[%s] --- Can't find comment from Request[%s]", itemID, r.Request.URL)
		return
	}
	commentCount := commentResponse.CommentsCount[0]
	// TODO: 价格转换
	var (
		errList []error
		tmpErr  error
	)
	commentCount.CommentCount, tmpErr = ParseChinesePrice(commentCount.CommentCountStr)
	errList = append(errList, tmpErr)
	commentCount.GoodCount, tmpErr = ParseChinesePrice(commentCount.GoodCountStr)
	errList = append(errList, tmpErr)
	commentCount.GeneralCount, tmpErr = ParseChinesePrice(commentCount.GeneralCountStr)
	errList = append(errList, tmpErr)
	commentCount.PoorCount, tmpErr = ParseChinesePrice(commentCount.PoorCountStr)
	errList = append(errList, tmpErr)
	pkg.GenerateErrorFromList(&err, errList)

	commentSummary := &storage.CommentSummary{
		ItemID:       itemID,
		AverageScore: commentCount.AverageScore,
		CommentCount: commentCount.CommentCount,
		GoodCountStr: commentCount.GoodCountStr,
		GoodCount:    commentCount.GoodCount,
		GoodRate:     commentCount.GoodRate,
		GeneralCount: commentCount.GeneralCount,
		GeneralRate:  commentCount.GeneralRate,
		PoorCount:    commentCount.PoorCount,
		PoorRate:     commentCount.PoorRate,
	}

	// 存储ItemID到ItemStorageMap
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
	log.GLogger.Infof("item[%s] --- item save status: %d", itemID, status)
}
