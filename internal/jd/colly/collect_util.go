package colly

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func GetItemIDFromUrl(url string) (itemID string, err error) {
	itemIDReg := regexp.MustCompile(`item\.jd\.com/(?P<item_id>\d*)\.html`)
	tmpList := itemIDReg.FindStringSubmatch(url)
	if len(tmpList) > 1 {
		itemID = tmpList[1]
	} else {
		err = fmt.Errorf("find item_id from url error, url: %s", url)
	}
	return
}

func ParseChinesePrice(chinesePrice string) (price int64, err error) {
	priceStr := strings.TrimRight(chinesePrice, "+")
	priceStr = strings.ReplaceAll(priceStr, "万", "0000")

	price, err = strconv.ParseInt(priceStr, 10, 64)
	return
}
