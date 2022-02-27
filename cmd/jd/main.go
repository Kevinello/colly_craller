package main

import (
	"kevinello.ltd/kevinello/collycrawler/internal/jd/collect"
	"kevinello.ltd/kevinello/collycrawler/internal/jd/storage"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg"
)

func main() {
	// start collect job
	go collect.CollectItemUrl("口红", collect.ItemUrlChan, 10)

	// ready to storage item
	go storage.StartStorageItem(storage.ItemStorageChan)

	// ready to collect item
	go collect.StartCollectItem(collect.ItemUrlChan)

	// 程序结束运行，关闭ticker
	defer close(pkg.StopChan)
	defer pkg.CloseTicker()
}
