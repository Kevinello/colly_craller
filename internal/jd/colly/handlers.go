package colly

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

const (
	priceSelectorFormat = "span.price.J-p-%s"
)

// HandlerFindItemIdFromUrl 从item url中获取item id
// @param r
// @author: Kevineluo
func HandlerFindItemIdFromUrl(r *colly.Request) {
	reg := regexp.MustCompile(`item\.jd\.com/(?P<item_id>\d*)\.html`)
	log.GLogger.Infof("visiting %s", r.URL.String())
	tmpList := reg.FindStringSubmatch(r.URL.String())
	if len(tmpList) > 1 {
		log.GLogger.Infof("find item_id: %s", tmpList[1])
		r.Ctx.Put("item_id", tmpList[1])
	} else {
		log.GLogger.Errorf("find item_id from url error, url: %s", r.URL.String())
	}
}

func HandlerParseItemInfo(h *colly.HTMLElement) {
	itemId := h.Request.Ctx.Get("item_id")
	priceSelector := fmt.Sprintf(priceSelectorFormat, itemId)
	log.GLogger.Infof("Get priceSelector: %s", priceSelector)
	price := h.DOM.Find(priceSelector).Text()
	log.GLogger.Infof("Price of item [%s]: %s", itemId, price)
}
