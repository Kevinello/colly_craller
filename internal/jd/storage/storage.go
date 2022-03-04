package storage

import (
	"encoding/json"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"kevinello.ltd/kevinello/collycrawler/internal/pkg/log"
)

var (
	DB *gorm.DB
)

func init() {
	// TODO: 初始化数据库
	dsn := "host=localhost user=kevinello password=jdCrawler2022 dbname=jd_data port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.GLogger.Alert("can't connect to database")
		os.Exit(1)
	}
}

func (item *Item) save() (err error) {
	itemStr, err := json.Marshal(item)
	log.GLogger.Info(string(itemStr))
	return
}
