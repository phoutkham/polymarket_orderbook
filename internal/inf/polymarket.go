package inf

type PriceLevel struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

type OrderBook struct {
	TokenId    string `json:"asset_id"`
	ServerTime string `json:"timestamp"`
	LocalTime  int64
	Bids       []PriceLevel `json:"bids"`
	Asks       []PriceLevel `json:"asks"`
}

type Market struct {
	Slug       string
	StartTime  int64
	EndTime    int64
	YesTokenId string
	NoTokenId  string
}

type Event struct {
	ConditionId   string  `json:"conditionId"`
	Outcomes      string  `json:"outcomes"`
	OutcomePrices string  `json:"outcomePrices"`
	ClobTokenIds  string  `json:"clobTokenIds"`
	StartTime     string  `json:"eventStartTime"`
	TakerBaseFee  int64   `json:"takerBaseFee"`
	MakerBaseFee  int64   `json:"makerBaseFee"`
	TickSize      float64 `json:"orderPriceMinTickSize"`
	NegRisk       bool    `json:"negRisk"`
}

type Token struct {
	TokenId string `json:"token_id"`
}

type Price struct {
	Price string `json:"price"`
}
