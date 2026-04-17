package helper

import (
	"encoding/binary"
	"math"
	"polymarket_orderbook/internal/inf"
	"strconv"
	"time"
)

func DumpOrderbook(o *inf.OrderBook, offset int64, path string) {
	var data [101]int32
	timestamp := make([]byte, len(data)*8)
	binary.LittleEndian.PutUint64(timestamp[:8], uint64(time.Now().UnixMilli()))

	for _, bid := range o.Bids {
		price, _ := strconv.ParseFloat(bid.Price, 64)
		size, _ := strconv.ParseFloat(bid.Size, 64)

		var index int32 = int32(math.Floor(price * 100))
		data[index] = data[index] + int32(size)
	}

	for _, ask := range o.Asks {
		price, _ := strconv.ParseFloat(ask.Price, 64)
		size, _ := strconv.ParseFloat(ask.Size, 64)

		var index int32 = int32(math.Ceil(price * 100))
		data[index] = data[index] - int32(size)
	}

	bytes := make([]byte, len(data)*4)

	for i, v := range data {
		binary.LittleEndian.PutUint32(bytes[i*4:], uint32(v))
	}

	orderBookPath := path + "/orderbook.bin"
	timestampPath := path + "/timestamp.bin"
	orderbookOffset := offset * (101 * 4)
	timestampOffset := offset * 8

	WriteToFile(orderbookOffset, &bytes, orderBookPath)
	WriteToFile(timestampOffset, &timestamp, timestampPath)

}
