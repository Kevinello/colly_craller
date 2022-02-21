package colly

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"kevinello.ltd/kevinello/collycrawler/internal/jd/storage"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

const (
	priceUrlFormatter   = "https://p.3.cn/prices/mgets?skuIds=J_%s"
	commentUrlFormatter = "https://club.jd.com/comment/productCommentSummaries.action?referenceIds=%s"
)

// HandlerFindItemIdFromUrl 从item url中获取item id
// @param r
// @author: Kevineluo
func HandlerFindItemIdFromUrl(r *colly.Request) {
	reg := regexp.MustCompile(`item\.jd\.com/(?P<item_id>\d*)\.html`)
	log.GLogger.Infof("visiting %s", r.URL.String())
	tmpList := reg.FindStringSubmatch(r.URL.String())
	if len(tmpList) > 1 {
		itemId := tmpList[1]
		log.GLogger.Infof("find item_id: %s", itemId)

		priceUrl := fmt.Sprintf(priceUrlFormatter, itemId)
		PriceCollector.Visit(priceUrl)

		commentUrl := fmt.Sprintf(commentUrlFormatter, itemId)
		CommentCollector.Visit(commentUrl)
	} else {
		log.GLogger.Errorf("find item_id from url error, url: %s", r.URL.String())
	}

}

// HandlerCollectPrice 从price接口收集商品价格
// @param r
// @author: Kevineluo
func HandlerCollectPrice(r *colly.Response) {
	// log.GLogger.Debugf("Get response: %s", string(r.Body))
	jsonStr := string(r.Body)
	jsonStr = strings.TrimSpace(jsonStr)

	priceResponse := storage.PriceResponse{}
	err := json.Unmarshal([]byte(jsonStr), &priceResponse)
	if err != nil {
		log.GLogger.Errorf("error when Unmarshal PriceResponse of Request[%s]: %s", r.Request.URL, err.Error())
		return
	}
	if len(priceResponse) > 0 {
		log.GLogger.Infof("Collect price from Request[%s]: %+v", r.Request.URL, priceResponse[0])
	} else {
		log.GLogger.Errorf("Can't find price from Request[%s]", r.Request.URL)
	}
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
	}
}
