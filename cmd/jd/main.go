package main

import (
	"kevinello.ltd/kevinello/collycrawler/internal/jd/web"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/anticrawl"
)

func main() {
	defer anticrawl.SeleniumService.Stop()
	web.CollectItemUrl("口红", web.ItemUrlChan)
}
