package main

import (
	"fmt"
	"log"

	"github.com/tsilvers/realtime-securities/markets/quote"

	"github.com/tsilvers/realtime-securities/markets/stock"
	"github.com/tsilvers/realtime-securities/provider"
)

const stocksPerRequest = 100

// main retrieves realtime quotes for requested stocks and displays them in ascending price order.
func main() {
	ds := provider.GetProvider("Tradier")

	// Create stocks for each symbol.
	var stocks = make(map[string]*stock.Stock)
	var symbols = make([]string, 0, 4)
	for _, symbol := range stock.GetSymbols() {
		symbols = append(symbols, symbol)
		stocks[symbol] = stock.NewStock(symbol)
	}

	if len(symbols) == 0 {
		log.Fatalln("No stocks found.")
	}

	// Retrieve quotes from the data provider and add them to the underlying stocks.
	for i := 0; i < len(symbols); i += stocksPerRequest {
		last := i + stocksPerRequest
		if last > len(stocks) {
			last = len(stocks)
		}

		quotes := ds.GetQuotes(symbols[i:last])

		for _, q := range quotes {
			stocks[q.Symbol()].SetQuote(q)
		}
	}

	// Sort stocks by price and display.
	stockList := make([]*stock.Stock, 0, len(stocks))
	for _, st := range stocks {
		stockList = append(stockList, st)
	}
	stock.SortByLastPrice(stockList)
	fmt.Print(quote.QuoteHeader())
	for _, st := range stockList {
		fmt.Println(st.Quote())
	}
}
