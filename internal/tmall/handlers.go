package tmall

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"kevinello.ltd/kevinello/collycrawler/internal/pkg/anticrawl"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
	"kevinello.ltd/kevinello/collycrawler/internal/tmall/storage"

	"github.com/gocolly/colly"
	"github.com/mitchellh/mapstructure"
)

func HandlerPrintRequestUrl(r *colly.Request) {
	log.GLogger.Infof("visiting %s", r.URL.String())
}

func HandlerSetCookie(r *colly.Request) {
	log.GLogger.Infof("visiting %s, will set cookie for the request", r.URL.String())
	cookieStr, err := anticrawl.GetCookieStr(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	if err != nil {
		log.GLogger.Errorf("get cookieStr error: %s", err.Error())
		return
	}
	r.Headers.Set("cookie", cookieStr)
}

func HandlerPrintResponseUrl(r *colly.Response) {
	log.GLogger.Infof("visited %s", r.Request.URL.String())
}

func HandlerPrintResponseContent(r *colly.Response) {
	log.GLogger.Infof("visited %s", r.Request.URL.String())
	log.GLogger.Infof("content: %s", string(r.Body))
}

func HandlerGetCategoryUrl(h *colly.HTMLElement) {
	log.GLogger.Info("script found")
	var err error

	jsonStr := strings.TrimPrefix(strings.TrimSpace(h.Text), "window.$data = ")
	tmpMap := make(map[string]interface{})

	if err = json.Unmarshal([]byte(jsonStr), &tmpMap); err != nil {
		log.GLogger.Error(err.Error())
		return
	}
	categoryMainLines, ok := tmpMap["2567006790"].(map[string]interface{})["categoryMainLines"]
	if !ok {
		err = fmt.Errorf("illegal script in the response of %s", h.Request.URL)
		log.GLogger.Error(err.Error())
		return
	}
	categories := make([]storage.Category, 0)
	if err = mapstructure.Decode(categoryMainLines, &categories); err != nil {
		log.GLogger.Error(err.Error())
		return
	}
	for _, category := range categories {
		if category.Action1 != "" {
			log.GLogger.Infof("found category [%s], will add url [%v] to CategoryCrawlQueue", category.Title1, category.Action1)
			CategoryCrawlQueue.AddURL(category.Action1)
		}
		if category.Action2 != "" {
			log.GLogger.Infof("found category [%s], will add url [%v] to CategoryCrawlQueue", category.Title2, category.Action2)
			CategoryCrawlQueue.AddURL(category.Action2)
		}
		if category.Action3 != "" {
			log.GLogger.Infof("found category [%s], will add url [%v] to CategoryCrawlQueue", category.Title3, category.Action3)
			CategoryCrawlQueue.AddURL(category.Action3)
		}
	}
	return
}
