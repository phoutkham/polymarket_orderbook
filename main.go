package main

import (
	"polymarket_orderbook/internal/polymarket"
	"polymarket_orderbook/internal/watcher"
	"strconv"
	"time"
)

func main() {
	var slugs []string = []string{
		"sea-sas-com-2026-04-17-com",
		"atp-sil-houkes-2026-04-17",
		"nba-cha-orl-2026-04-17",
		"nba-gsw-phx-2026-04-17",
		"cs2-vit-navi-2026-04-17",
		"dota2-ts8-vg-2026-04-18",
		"lol-ae2-tts1-2026-04-17",
		"val-c9-nv2-2026-04-17",
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

		time.Sleep(500 * time.Millisecond)
	}

}
