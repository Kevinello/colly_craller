package storage

type Item struct {
	ItemID string
	Price  int64
	SkuNum int64
}

type PriceResponse []struct {
	Price string `json:"p"`
	Op    string `json:"op"`
	Cbf   string `json:"cbf"`
	ID    string `json:"id"`
	M     string `json:"m"`
}

type CommentResponse struct {
	CommentsCount []struct {
		SkuID               int     `json:"SkuId"`
		ProductID           int     `json:"ProductId"`
		ShowCount           int     `json:"ShowCount"`
		ShowCountStr        string  `json:"ShowCountStr"`
		CommentCountStr     string  `json:"CommentCountStr"`
		CommentCount        int     `json:"CommentCount"`
		AverageScore        int     `json:"AverageScore"`
		DefaultGoodCountStr string  `json:"DefaultGoodCountStr"`
		DefaultGoodCount    int     `json:"DefaultGoodCount"`
		GoodCountStr        string  `json:"GoodCountStr"`
		GoodCount           int     `json:"GoodCount"`
		AfterCount          int     `json:"AfterCount"`
		OneYear             int     `json:"OneYear"`
		AfterCountStr       string  `json:"AfterCountStr"`
		VideoCount          int     `json:"VideoCount"`
		VideoCountStr       string  `json:"VideoCountStr"`
		GoodRate            float64 `json:"GoodRate"`
		GoodRateShow        int     `json:"GoodRateShow"`
		GoodRateStyle       int     `json:"GoodRateStyle"`
		GeneralCountStr     string  `json:"GeneralCountStr"`
		GeneralCount        int     `json:"GeneralCount"`
		GeneralRate         float64 `json:"GeneralRate"`
		GeneralRateShow     int     `json:"GeneralRateShow"`
		GeneralRateStyle    int     `json:"GeneralRateStyle"`
		PoorCountStr        string  `json:"PoorCountStr"`
		PoorCount           int     `json:"PoorCount"`
		SensitiveBook       int     `json:"SensitiveBook"`
		PoorRate            float64 `json:"PoorRate"`
		PoorRateShow        int     `json:"PoorRateShow"`
		PoorRateStyle       int     `json:"PoorRateStyle"`
	} `json:"CommentsCount"`
}
