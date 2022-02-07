package colly

import (
	"regexp"

	"github.com/gocolly/colly"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

func HandlerFindItemIdFromUrl(r *colly.Request) {
	reg := regexp.MustCompile(`^item\.jd\.com/(?P<item_id>\d{12})\.html`)
	log.GLogger.Infof("visiting %s", r.URL.String())
	tmpList := reg.FindStringSubmatch(r.URL.String())
	if len(tmpList) > 1 {
		log.GLogger.Infof("find item_id: %s", tmpList[1])
		r.Ctx.Put("item_id", tmpList[1])
	} else {
		log.GLogger.Errorf("find item_id from url error, url: %s", r.URL.String())
	}
}
