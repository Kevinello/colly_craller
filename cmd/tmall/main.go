package main

import (
	"kevinello.ltd/kevinello/collycrawller/internal/pkg/log"
	"kevinello.ltd/kevinello/collycrawller/internal/tmall"
)

func main() {
	// Start scraping on https://www.tmall.com
	err := tmall.HomePageCollector.Visit("https://www.tmall.com/")
	if err != nil {
		log.GLogger.Error(err.Error())
	}
	log.GLogger.Info("Crawlling job completed")
}
