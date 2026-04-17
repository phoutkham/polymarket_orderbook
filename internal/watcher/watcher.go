package watcher

import (
	"polymarket_orderbook/internal/helper"
	"polymarket_orderbook/internal/inf"
	"polymarket_orderbook/internal/polymarket"
	"time"
)

type List struct {
	Market    *inf.Market
	StartTime int64
	EndTime   int64
	Offset    int64
	IsClosed  bool
	Callback  func(*inf.OrderBook, int64, string)
}

type Watcher struct {
	Lists []*List
}

func NewWatcher() *Watcher {
	return &Watcher{
		Lists: []*List{},
	}
}

func (w *Watcher) Add(slug string, p *polymarket.Polymarket) {

	market, err := p.GetMarket(slug)

	if err != nil {
		return
	}

	w.Lists = append(w.Lists, &List{
		Market:    market,
		StartTime: market.StartTime,
		//StartTime: 0,
		EndTime: market.StartTime + (60 * 60 * 3),
		//EndTime:  100000000000000,
		IsClosed: false,
	})

}

func (w *Watcher) GetAllTokenToWatch() *[]inf.Token {
	var tokenIds []inf.Token
	var currentTime int64 = time.Now().Unix()

	for _, list := range w.Lists {

		if list.Offset > 100000 {
			continue
		}

		if currentTime > list.StartTime && currentTime < list.EndTime && !list.IsClosed {
			tokenId := inf.Token{TokenId: list.Market.YesTokenId}
			tokenIds = append(tokenIds, tokenId)
		}

	}

	return &tokenIds
}

func (l *List) Save(o *inf.OrderBook) {
	var path string = "./data/" + l.Market.Slug

	helper.DumpOrderbook(o, l.Offset, path)

	l.Offset = l.Offset + 1
}
