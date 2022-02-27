package storage

import (
	"encoding/json"

	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

var ()

func init() {
	// TODO: 初始化数据库
}

func (item *Item) save() (err error) {
	itemStr, err := json.Marshal(item)
	log.GLogger.Info(string(itemStr))
	return
}
