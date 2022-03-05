package main

import (
	"kevinello.ltd/kevinello/collycrawler/internal/jd/collect"
	"kevinello.ltd/kevinello/collycrawler/internal/jd/storage"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg"
)

func main() {
	// ready to storage item
	go storage.StartStorageItem(storage.ItemStorageChan)
	// ready to collect item
	go collect.StartCollectItem(collect.ItemUrlChan)
	// start collect job
	go collect.CollectItemUrl(pkg.GetEnv("SEARCH_KEYWORD", "美妆"), collect.ItemUrlChan, 100)

	// 程序结束运行，关闭ticker
	defer close(pkg.StopChan)
	defer pkg.CloseTicker()
}
