package colly

import (
	"fmt"
	"regexp"
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
