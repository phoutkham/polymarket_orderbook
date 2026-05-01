package main

import (
	"polymarket_orderbook/internal/polymarket"
	"polymarket_orderbook/internal/watcher"
	"time"
)

func main() {

	polymarket := polymarket.NewPolymarketAPI()
	watcher := watcher.NewWatcher(polymarket)

	var slugs []string = []string{
		polymarket.FetchCurrentSlug(),
		polymarket.FetchNextSlug(),
	}

	for _, s := range slugs {
		watcher.Add(s, polymarket)
	}

	for {
		o, err := polymarket.GetOrderBooks(watcher.GetAllTokenToWatch())

		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		for _, o := range *o {
			for _, l := range watcher.Lists {
				if l.Market.YesTokenId == o.TokenId {
					l.Save(&o)
				}
			}
		}

		watcher.Clean(polymarket)

		time.Sleep(100 * time.Millisecond)
	}

}
