package collycrawller

import "github.com/gocolly/colly"

func HandlerPrintUrl(r *colly.Request) {
	GLogger.Infof("visiting %s", r.URL.String())
}
