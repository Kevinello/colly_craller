package colly

import (
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"

	"github.com/gocolly/colly"
)

func HandlerPrintRequestUrl(r *colly.Request) {
	log.GLogger.Infof("visiting %s", r.URL.String())
}

func HandlerPrintResponseUrl(r *colly.Response) {
	log.GLogger.Infof("visited %s", r.Request.URL.String())
}

func HandlerPrintResponseContent(r *colly.Response) {
	log.GLogger.Infof("visited %s", r.Request.URL.String())
	log.GLogger.Infof("content: %s", string(r.Body))
}
