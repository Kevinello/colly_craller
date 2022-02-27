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
	Price          int64
	SkuNum         int64
	CommentSummary *CommentSummary
}

type CommentSummary struct {
	AverageScore int64
	CommentCount int64
	GoodCount    int64
	GoodRate     float64
	GeneralCount int64
	GeneralRate  float64
	PoorCount    int64
	PoorRate     float64
}

// ItemSaveMessage Item存储
type ItemSaveMessage struct {
	ItemID    string
	SaveField string
	SaveValue interface{}
	SaveRes   chan int
}

// PriceResponse 价格查询接口返回结构体
type PriceResponse []struct {
	Price string `json:"p"`
	Op    string `json:"op"`
	Cbf   string `json:"cbf"`
	ID    string `json:"id"`
	M     string `json:"m"`
}

// CommentResponse 评论查询接口返回结构体
type CommentResponse struct {
	CommentsCount []struct {
		SkuID               int64   `json:"SkuId"`
		ProductID           int64   `json:"ProductId"`
		ShowCount           int64   `json:"ShowCount"`
		ShowCountStr        string  `json:"ShowCountStr"`
		CommentCountStr     string  `json:"CommentCountStr"`
		CommentCount        int64   `json:"CommentCount"`
		AverageScore        int64   `json:"AverageScore"`
		DefaultGoodCountStr string  `json:"DefaultGoodCountStr"`
		DefaultGoodCount    int64   `json:"DefaultGoodCount"`
		GoodCountStr        string  `json:"GoodCountStr"`
		GoodCount           int64   `json:"GoodCount"`
		AfterCount          int64   `json:"AfterCount"`
		OneYear             int64   `json:"OneYear"`
		AfterCountStr       string  `json:"AfterCountStr"`
		VideoCount          int64   `json:"VideoCount"`
		VideoCountStr       string  `json:"VideoCountStr"`
		GoodRate            float64 `json:"GoodRate"`
		GoodRateShow        int64   `json:"GoodRateShow"`
		GoodRateStyle       int64   `json:"GoodRateStyle"`
		GeneralCountStr     string  `json:"GeneralCountStr"`
		GeneralCount        int64   `json:"GeneralCount"`
		GeneralRate         float64 `json:"GeneralRate"`
		GeneralRateShow     int64   `json:"GeneralRateShow"`
		GeneralRateStyle    int64   `json:"GeneralRateStyle"`
		PoorCountStr        string  `json:"PoorCountStr"`
		PoorCount           int64   `json:"PoorCount"`
		SensitiveBook       int64   `json:"SensitiveBook"`
		PoorRate            float64 `json:"PoorRate"`
		PoorRateShow        int64   `json:"PoorRateShow"`
		PoorRateStyle       int64   `json:"PoorRateStyle"`
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
					log.GLogger.Errorf("error when saveField for item[%s], error: %s", message.ItemID, err.Error())
					// 返回存储状态
					message.SaveRes <- StatusItemSaveFieldError
					continue
				} else {
					log.GLogger.Infof("saveField for item[%s] success", message.ItemID)
					itemStorage.savedFieldSet.Add(message.SaveField)
					// 若存储字段数达到Item字段数，则进行持久化
					if itemStorage.savedFieldSet.Cardinality() >= reflect.TypeOf(itemStorage.item).Elem().NumField() {
						log.GLogger.Infof("item[%s] is ready(savedFieldNum: %d)", message.ItemID, itemStorage.savedFieldSet.Cardinality())
						err := itemStorage.item.save()
						if err != nil {
							log.GLogger.Errorf("error when save item[%s], error: %s", message.ItemID, err.Error())
							// 返回存储状态
							message.SaveRes <- StatusItemStorageError
							continue
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
