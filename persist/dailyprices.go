package persist

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/tsilvers/realtime-securities/markets/stock"
	"io"
	"os"
)

func InitPricesStore() {
	_ = os.RemoveAll(pricesDir)
	_ = os.MkdirAll(pricesDir, 0755)
}

func SavePrices(symbol string, prices []stock.DailyPriceGob) error {
	fname := pricesDir + stockFilename(symbol)
	file, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("could not persist prices for %s to file %s: %w", symbol, fname, err)
	}
	defer file.Close()

	enc := gob.NewEncoder(file)

	for _, price := range prices {
		if err := enc.Encode(price); err != nil {
			return fmt.Errorf("could not persist prices for %s to file %s: %w", symbol, fname, err)
		}
	}

	return nil
}

func LoadPrices(symbol string) ([]*stock.DailyPriceGob, error) {
	var prices []*stock.DailyPriceGob

	fname := pricesDir + stockFilename(symbol)
	file, err := os.Open(fname)
	if err != nil {
		return prices, fmt.Errorf("could not load prices for %s from file %s: %w", symbol, fname, err)
	}
	defer file.Close()

	enc := gob.NewDecoder(file)

	for err == nil {
		price := &stock.DailyPriceGob{}
		if err = enc.Decode(price); err == nil {
			prices = append(prices, price)
		}
	}

	if !errors.Is(io.EOF, err) {
		return nil, fmt.Errorf("could not load prices for %s from file %s: %w", symbol, fname, err)
	}

	return prices, nil
}
