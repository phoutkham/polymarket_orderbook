package polymarket

import (
	"encoding/json"
	"polymarket_orderbook/internal/http"
	"polymarket_orderbook/internal/inf"
	"strconv"
	"time"
)

type Polymarket struct {
	Http *http.Client
}

func NewPolymarketAPI() *Polymarket {
	return &Polymarket{
		Http: http.NewClient(),
	}
}

func (p *Polymarket) FetchCurrentSlug() string {

	_time := (time.Now().Unix() / 300) * 300
	timeString := strconv.FormatInt(_time, 10)

	return "btc-updown-5m-" + timeString

}

func (p *Polymarket) FetchNextSlug() string {

	_time := ((time.Now().Unix() / 300) * 300) + 300
	timeString := strconv.FormatInt(_time, 10)

	return "btc-updown-5m-" + timeString

}

func (p *Polymarket) GetPrice() *int64 {

	var price inf.Price

	bytes, err := p.Http.GET("https://api.exchange.coinbase.com/products/BTC-USD/ticker")

	if err != nil {
		return nil
	}

	err = json.Unmarshal(bytes, &price)

	if err != nil {
		return nil
	}

	_price, _ := strconv.ParseFloat(price.Price, 64)
	var __price int64 = int64(_price)

	return &__price

}
func (p *Polymarket) GetOrderBook(tokenId string) (*inf.OrderBook, error) {
	var _orderbook inf.OrderBook

	bytes, err := p.Http.GET("https://clob.polymarket.com/book?token_id=" + tokenId)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &_orderbook)

	if err != nil {
		return nil, err
	}

	_orderbook.LocalTime = time.Now().UnixMilli()

	return &_orderbook, nil

}

func (p *Polymarket) GetMarket(slug string) (*inf.Market, error) {
	var event []inf.Event

	bytes, err := p.Http.GET("https://gamma-api.polymarket.com/markets?slug=" + slug + "&closed=false")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &event)

	if err != nil {
		return nil, err
	}

	if len(event) == 0 {
		return nil, nil
	}

	startTime, err := time.Parse(time.RFC3339, event[0].StartTime)

	var clobTokenIds []string

	if err := json.Unmarshal([]byte(event[0].ClobTokenIds), &clobTokenIds); err != nil {
		return nil, err
	}

	return &inf.Market{
		Slug:       slug,
		StartTime:  startTime.Unix(),
		EndTime:    0,
		YesTokenId: clobTokenIds[0],
		NoTokenId:  clobTokenIds[1],
	}, nil

}

func (p *Polymarket) GetOrderBooks(tokenIds interface{}) (*[]inf.OrderBook, error) {
	var _orderbook []inf.OrderBook

	bytes, err := p.Http.POST("https://clob.polymarket.com/books", tokenIds)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &_orderbook)

	if err != nil {
		return nil, err
	}

	return &_orderbook, nil

}
