package main

import (
	"fmt"
	"github.com/tsilvers/realtime-securities/markets/stock"
	"github.com/tsilvers/realtime-securities/persist"
	"github.com/tsilvers/realtime-securities/provider"
	"log"
	"os"
	"time"
)

// main retrieves and persists daily stock prices starting from the requested date.
func main() {
	ds := provider.GetProvider("Tradier")

	// Get start date.
	if len(os.Args) != 2 {
		usage()
	}
	start, err := time.Parse("01/02/2006", os.Args[1])
	if err != nil {
		usage()
	}

	fmt.Println("Loading daily stock prices...")

	// Retrieve list of stock symbols.
	symbols := stock.GetSymbols()

	// Clear stored prices.
	persist.InitPricesStore()

	// Load price histories.
	var prices []*stock.DailyPrice
	cnt := 0
	for _, symbol := range symbols {
		cnt++
		fmt.Printf("%4d: %s\n", cnt, symbol)

		// Load price history from data source.
		prices = ds.GetPriceHistory(symbol, start)

		// Persist price history.
		var priceGobs []stock.DailyPriceGob
		for _, price := range prices {
			priceGobs = append(priceGobs, price.ToGob())
			fmt.Printf("      %s\n", price.Date().Format("01/02/2006"))
		}
		err = persist.SavePrices(symbol, priceGobs)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage: dailyprices StartDate (mm/dd/yyyy)\n\n")
	os.Exit(1)
}
