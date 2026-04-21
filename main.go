package main

import (
	"polymarket_orderbook/internal/polymarket"
	"polymarket_orderbook/internal/watcher"
	"strconv"
	"time"
)

func main() {
	var slugs []string = []string{
		"lal-mal-val-2026-04-21-mal",
		"lal-mal-val-2026-04-21-val",
		"lal-bil-osa-2026-04-21-bil",
		"lal-bil-osa-2026-04-21-osa",
	}

	polymarket := polymarket.NewPolymarketAPI()
	watcher := watcher.NewWatcher()

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

					bestBid, _ := strconv.ParseFloat(o.Bids[len(o.Bids)-1].Price, 10)

					if bestBid <= 0.01 || bestBid >= 0.99 {
						l.IsClosed = true
					}

				}
			}
		}

		time.Sleep(200 * time.Millisecond)
	}

}
