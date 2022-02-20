package main

import (
	"kevinello.ltd/kevinello/collycrawler/internal/jd/colly"
	"kevinello.ltd/kevinello/collycrawler/internal/jd/web"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/anticrawl"
)

func main() {

	// 最后关闭SeleniumService
	defer anticrawl.SeleniumService.Stop()
	go web.CollectItemUrl("口红", web.ItemUrlChan)
	colly.ItemQueue.Run(colly.ItemCollector)
}
