package main

import (
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
	"kevinello.ltd/kevinello/collycrawler/internal/tmall"
)

func main() {
	// Start scraping on https://www.tmall.com
	err := tmall.HomePageCollector.Visit("https://www.tmall.com/")
	if err != nil {
		log.GLogger.Error(err.Error())
	}
	log.GLogger.Info("Crawlling job completed")
}
