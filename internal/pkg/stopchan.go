package pkg

import (
	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

var (
	// 无缓存只写chan，用于阻塞主进程
	StopChan = make(chan<- struct{})
)

// CloseTicker defer CloseTicker()阻塞主进程，接收到停止信号时，关闭ticker，结束主进程
func CloseTicker() {

	log.GLogger.Info("waiting for stop signal...")
	// 阻塞直到StopChan被close（因为StopChan无缓存只写，空结构体无法被接收）
	StopChan <- struct{}{}
	log.GLogger.Alert("receive stop signal, will stop programme")
}
