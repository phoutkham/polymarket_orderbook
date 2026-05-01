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
	BTC       *int64
	Callback  func(*inf.OrderBook, int64, string)
}

type Watcher struct {
	Lists []*List
	BTC   int64
}

func NewWatcher(p *polymarket.Polymarket) *Watcher {

	watcher := &Watcher{
		Lists: []*List{},
	}

	go func() {
		for {

			btc := p.GetPrice()
			watcher.BTC = *btc

			time.Sleep(50 * time.Millisecond)

		}
	}()

	return watcher
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
		EndTime: market.StartTime + (60 * 5),
		//EndTime:  100000000000000,
		IsClosed: false,
		BTC:      &w.BTC,
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
		} else if currentTime > list.EndTime && !list.IsClosed {
			list.IsClosed = true
		}

	}

	return &tokenIds
}

func (l *List) Save(o *inf.OrderBook) {
	var path string = "./data/" + l.Market.Slug
	var btc int64 = *l.BTC

	helper.DumpOrderbook(o, l.Offset, path, btc)

	l.Offset = l.Offset + 1
}

func (w *Watcher) Clean(p *polymarket.Polymarket) {
	var newList []*List

	var createNewSlug bool = false

	for _, list := range w.Lists {
		if list.IsClosed {
			createNewSlug = true
			continue
		}
		newList = append(newList, list)
	}

	w.Lists = newList

	if createNewSlug {
		w.Add(p.FetchNextSlug(), p)
	}

}
