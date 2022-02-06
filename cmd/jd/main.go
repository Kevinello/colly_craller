package jd

import "kevinello.ltd/kevinello/collycrawler/internal/pkg/anticrawl"

func main() {
	defer anticrawl.SeleniumService.Stop()
}
