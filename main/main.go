package main

import "collycrawller"

func main() {
	// Start scraping on https://www.tmall.com
	err := collycrawller.HomePageCollector.Visit("https://www.tmall.com/")
	if err != nil {
		collycrawller.GLogger.Error(err.Error())
	}
	collycrawller.GLogger.Info("Crawlling job completed")
}
