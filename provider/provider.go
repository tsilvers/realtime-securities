package provider

import (
	"github.com/tsilvers/realtime-securities/markets"
	"github.com/tsilvers/realtime-securities/markets/option"
	"github.com/tsilvers/realtime-securities/markets/quote"
	"github.com/tsilvers/realtime-securities/markets/stock"
	"time"

	"github.com/tsilvers/realtime-securities/provider/tradier"
)

type Provider interface {
	GetPriceHistory(symbol string, start time.Time) []*stock.DailyPrice
	GetTimeSales(symbol string, start time.Time) []stock.OneMinSale
	GetQuotes(symbols []string) []quote.Quote
	GetOptionExpirations(symbol string) []option.Expiration
	GetMarketStatus() markets.StatusType
}

var providers = map[string]Provider{
	"Tradier": tradier.New(),
}

func GetProvider(name string) Provider {
	return providers[name]
}
