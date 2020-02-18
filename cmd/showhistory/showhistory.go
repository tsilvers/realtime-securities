package main

import (
	"fmt"
	"log"

	"github.com/tsilvers/realtime-securities/markets/stock"
	"github.com/tsilvers/realtime-securities/persist"
)

// main retrieves and displays persisted daily stock price info.
func main() {
	fmt.Println("Retrieving daily stock prices...")

	// Retrieve list of stock symbols.
	symbols := stock.GetSymbols()

	// Create stocks and load price history from persistence.
	var prices []*stock.DailyPrice
	var stocks = make(map[string]*stock.Stock)
	for _, symbol := range symbols {
		stocks[symbol] = stock.NewStock(symbol)

		priceGobs, err := persist.LoadPrices(symbol)
		if err != nil {
			log.Printf("Persistence load error %s: %s", symbol, err)
			continue
		}

		prices = nil
		for _, priceGob := range priceGobs {
			price, err := priceGob.ToDailyPrice()
			if err != nil {
				log.Printf("Error loading price history for %s: %s\n", symbol, err)
				continue
			}
			prices = append(prices, &price)
		}

		if err = stocks[symbol].AppendDailyPrices(prices); err != nil {
			log.Printf("Persistence load error %s: %s", symbol, err)
			continue
		}
	}

	// Display daily prices.
	stockList := make([]*stock.Stock, 0, len(stocks))
	for _, st := range stocks {
		stockList = append(stockList, st)
	}
	stock.SortBySymbol(stockList)
	for _, st := range stockList {
		fmt.Printf("\n%s:\n", st.Symbol())

		fmt.Print(stock.DailyPriceHeader())
		for _, dp := range st.Prices() {
			fmt.Println(dp)
		}
	}
}
