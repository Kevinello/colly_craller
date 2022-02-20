package main

import (
	"kevinello.ltd/kevinello/collycrawler/internal/jd/colly"
	"kevinello.ltd/kevinello/collycrawler/internal/jd/web"
	"kevinello.ltd/kevinello/collycrawler/internal/pkg"
)

func main() {
	go web.CollectItemUrl("口红", web.ItemUrlChan, 10)
	// TODO: 需要让Queue不直接结束
	colly.ItemQueue.Run(colly.ItemCollector)

	// 程序结束运行，关闭ticker
	defer close(pkg.StopChan)
	defer pkg.CloseTicker()
}
