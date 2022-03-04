package storage

import (
	"fmt"
	"reflect"

	mapset "github.com/deckarep/golang-set"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

type ItemStorage struct {
	item          *Item
	savedFieldSet mapset.Set
}

type Item struct {
	ItemID         string
	SkuNum         int
	Price          *Price
	ShopInfo       *ShopInfo
	Promotion      *Promotion
	CommentSummary *CommentSummary
}

type Price struct {
	HagglePromotion bool   `json:"hagglePromotion"`
	Op              string `json:"op"`
	P               string `json:"p"`
}

type ShopInfo struct {
	Shop struct {
		AvgEfficiencyScore   int    `json:"avgEfficiencyScore"`
		AvgServiceScore      int    `json:"avgServiceScore"`
		AvgWareScore         int    `json:"avgWareScore"`
		CateGoodShop         int    `json:"cateGoodShop"`
		Diamond              bool   `json:"diamond"`
		FollowCount          int    `json:"followCount"`
		GoodShop             int    `json:"goodShop"`
		HasCoupon            bool   `json:"hasCoupon"`
		Name                 string `json:"name"`
		NewNum               int    `json:"newNum"`
		PromotionNum         int    `json:"promotionNum"`
		Score                int    `json:"score"`
		ServiceScore         int    `json:"serviceScore"`
		ShopActivityTotalNum int    `json:"shopActivityTotalNum"`
		ShopID               int    `json:"shopId"`
		TotalNum             int    `json:"totalNum"`
		WareScore            int    `json:"wareScore"`
	} `json:"shop"`
	CustomerService struct {
		HasChat bool `json:"hasChat"`
		HasJimi bool `json:"hasJimi"`
		Online  bool `json:"online"`
	} `json:"customerService"`
}

type Promotion struct {
	Activity []struct {
		ActivityType string `json:"activityType"`
		CustomTag    struct {
		} `json:"customTag"`
		PromoID    string `json:"promoId"`
		Text       string `json:"text"`
		TypeNumber string `json:"typeNumber"`
		Value      string `json:"value"`
	} `json:"activity"`
}

type CommentSummary struct {
	AverageScore    int64
	CommentCountStr string
	CommentCount    int64
	GoodCountStr    string
	GoodCount       int64
	GoodRate        float64
	GeneralCountStr string
	GeneralCount    int64
	GeneralRate     float64
	PoorCountStr    string
	PoorCount       int64
	PoorRate        float64
}

// ItemSaveMessage Item存储
type ItemSaveMessage struct {
	ItemID    string
	SaveField string
	SaveValue interface{}
	SaveRes   chan int
}

// WareBussinessResponse 商品综合信息查询接口返回结构体
type WareBussinessResponse struct {
	RankUnited struct {
		RevertItem struct {
			ID          string `json:"id"`
			Jump        string `json:"jump"`
			JumpTypeInt int    `json:"jumpTypeInt"`
			Name        string `json:"name"`
			RankID      string `json:"rankId"`
			RankTypeInt int    `json:"rankTypeInt"`
		} `json:"revertItem"`
	} `json:"rankUnited"`
	Price     Price     `json:"price"`
	ShopInfo  ShopInfo  `json:"shopInfo"`
	Promotion Promotion `json:"promotion"`
}

// CommentResponse 评论查询接口返回结构体
type CommentResponse struct {
	CommentsCount []struct {
		ShowCount           int64   `json:"ShowCount"`
		ShowCountStr        string  `json:"ShowCountStr"`
		CommentCountStr     string  `json:"CommentCountStr"`
		CommentCount        int64   `json:"CommentCount"`
		AverageScore        int64   `json:"AverageScore"`
		DefaultGoodCountStr string  `json:"DefaultGoodCountStr"`
		DefaultGoodCount    int64   `json:"DefaultGoodCount"`
		GoodCountStr        string  `json:"GoodCountStr"`
		GoodCount           int64   `json:"GoodCount"`
		GoodRate            float64 `json:"GoodRate"`
		AfterCount          int64   `json:"AfterCount"`
		OneYear             int64   `json:"OneYear"`
		AfterCountStr       string  `json:"AfterCountStr"`
		VideoCount          int64   `json:"VideoCount"`
		VideoCountStr       string  `json:"VideoCountStr"`
		GeneralCountStr     string  `json:"GeneralCountStr"`
		GeneralCount        int64   `json:"GeneralCount"`
		GeneralRate         float64 `json:"GeneralRate"`
		PoorCountStr        string  `json:"PoorCountStr"`
		PoorCount           int64   `json:"PoorCount"`
		PoorRate            float64 `json:"PoorRate"`
		PoorRateShow        int64   `json:"PoorRateShow"`
	} `json:"CommentsCount"`
}

var (
	ItemStorageMap = make(map[string]*ItemStorage)
	// ItemStorageChan Item存储channel，保证ItemStorageMap的线程安全
	ItemStorageChan = make(chan *ItemSaveMessage, 30)
)

func init() {

}

// StartStorageItem 开始存储Item
// @param isc
// @author: Kevineluo
func StartStorageItem(isc chan *ItemSaveMessage) {
	for {
		select {
		case message := <-isc:
			log.GLogger.Debugf("receive item save message: %+v", *message)
			if message.ItemID == "" {
				log.GLogger.Errorf("illegal item save message, itemID not found, received message: %+v", message)
				// 返回存储状态
				message.SaveRes <- StatusItemIDNotFound
				continue
			}
			itemStorage, exist := ItemStorageMap[message.ItemID]
			// ItemStorageMap中不存在则初始化Item
			if !exist {
				fieldSet := mapset.NewSet()
				itemStorage = &ItemStorage{
					item:          &Item{ItemID: message.ItemID},
					savedFieldSet: fieldSet,
				}
				ItemStorageMap[message.ItemID] = itemStorage
			}
			// 存储字段
			if message.SaveField != "" && message.SaveValue != nil {
				err := saveField(itemStorage.item, message.SaveField, message.SaveValue)
				if err != nil {
					log.GLogger.Errorf("item[%s] --- error when saveField, error: %s", message.ItemID, err.Error())
					// 返回存储状态
					message.SaveRes <- StatusItemSaveFieldError
					continue
				} else {
					log.GLogger.Infof("saveField for item[%s] success", message.ItemID)
					itemStorage.savedFieldSet.Add(message.SaveField)
					// 若存储字段数达到Item字段数，则进行持久化
					log.GLogger.Debugf("item[%s] --- saved field num: %d", message.ItemID, itemStorage.savedFieldSet.Cardinality())
					if itemStorage.savedFieldSet.Cardinality() >= reflect.TypeOf(itemStorage.item).Elem().NumField()-1 {
						log.GLogger.Infof("item[%s] --- item is ready(savedFieldNum: %d)", message.ItemID, itemStorage.savedFieldSet.Cardinality())
						err := itemStorage.item.save()
						if err != nil {
							log.GLogger.Errorf("item[%s] --- error when save item, error: %s", message.ItemID, err.Error())
							// 返回存储状态
							message.SaveRes <- StatusItemStorageError
							continue
						} else {
							// item持久化成功，清理ItemStorageMap
							delete(ItemStorageMap, message.ItemID)
						}
					}
				}
			}
			// 返回存储状态
			message.SaveRes <- StatusItemSavedSuccess
		}
	}
}

// saveField 存储Item的一个字段
// @param item
// @param field
// @param value
// @return err
// @author: Kevineluo
func saveField(item *Item, setField string, setValue interface{}) (err error) {
	val := reflect.ValueOf(item).Elem()
	toSetField := val.FieldByName(setField)
	if toSetField.IsValid() {
		toSetField.Set(reflect.ValueOf(setValue))
	} else {
		err = fmt.Errorf("can't find valid field[%s]", setField)
	}
	return
}
