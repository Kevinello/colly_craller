package main

import (
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/anticrawl"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
	"kevinello.ltd/kevinello/collycrawler/internal/tmall/colly"
)

func main() {
	// Start scraping on https://www.tmall.com
	defer anticrawl.SeleniumService.Stop()

	err := colly.HomePageCollector.Visit("https://www.tmall.com/")
	if err != nil {
		log.GLogger.Error(err.Error())
	}
	log.GLogger.Info("Crawlling job completed")
}
