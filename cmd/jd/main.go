package main

import (
	"kevinello.ltd/kevinello/collycrawler/internal/jd/collect"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg"
)

func main() {
	// start collect job
	go collect.CollectItemUrl("口红", collect.ItemUrlChan, 10)

	go collect.HandleItemUrl(collect.ItemUrlChan)
	// 程序结束运行，关闭ticker
	defer close(pkg.StopChan)
	defer pkg.CloseTicker()
}
