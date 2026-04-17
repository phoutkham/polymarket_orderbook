package polymarket

import (
	"encoding/json"
	"errors"
	"polymarket_orderbook/internal/http"
	"polymarket_orderbook/internal/inf"
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

	layout := "2006-01-02 15:04:05-07"
	startTime, _ := time.Parse(layout, event[0].StartTime)

	if startTime.Unix() < 0 {
		return nil, errors.New("error parsing startTime")
	}

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
