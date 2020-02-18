package stock

import (
	"fmt"
	"sort"
	"sync"

	"github.com/tsilvers/realtime-securities/markets/option"
	"github.com/tsilvers/realtime-securities/markets/quote"
)

type Stock struct {
	sync.Mutex
	symbol      string
	quote       quote.Quote
	optionChain []option.Expiration
	dailyPrices []DailyPrice
	oneMinSales []OneMinSale
}

func NewStock(symbol string) *Stock {
	return &Stock{
		symbol:      symbol,
		optionChain: make([]option.Expiration, 0, 16),
		dailyPrices: make([]DailyPrice, 0, 128),
		oneMinSales: make([]OneMinSale, 0, 128),
	}
}

func (s *Stock) Symbol() string {
	return s.symbol
}

func (s *Stock) Quote() quote.Quote {
	return s.quote
}

func (s *Stock) LastPrice() float64 {
	return s.quote.Last()
}

func (s *Stock) Prices() []DailyPrice {
	s.Lock()
	defer s.Unlock()

	dps := make([]DailyPrice, 0, len(s.dailyPrices))
	for _, dp := range s.dailyPrices {
		dps = append(dps, dp)
	}

	return dps
}

func (s *Stock) SetQuote(q quote.Quote) {
	s.Lock()
	defer s.Unlock()

	s.quote = q
}

func (s *Stock) AppendExpirations(exps []option.Expiration) error {
	s.Lock()
	defer s.Unlock()

	for _, exp := range exps {
		last := len(s.optionChain)
		if last > 0 {
			if !exp.Date().After(s.optionChain[last-1].Date()) {
				return fmt.Errorf("option expiration dates must be appended in chronological order - last date = %s, new date = %s",
					s.optionChain[last-1].ExpirationDateStr(), exp.ExpirationDateStr())
			}
		}

		s.optionChain = append(s.optionChain, exp)
	}

	return nil
}

func (s *Stock) AppendDailyPrices(dps []*DailyPrice) error {
	s.Lock()
	defer s.Unlock()

	for _, dp := range dps {

		last := len(s.dailyPrices)
		if last > 0 {
			if !dp.Date().After(s.dailyPrices[last-1].Date()) {
				return fmt.Errorf("daily prices must be appended in chronological order - last date = %s, new date = %s",
					s.dailyPrices[last-1].DateStr(), dp.DateStr())
			}
		}

		s.dailyPrices = append(s.dailyPrices, *dp)
	}

	return nil
}

func (s *Stock) AppendOneMinSales(omss []OneMinSale) error {
	s.Lock()
	defer s.Unlock()

	for _, oms := range omss {
		last := len(s.oneMinSales)
		if last > 0 {
			if !oms.StartTime().After(s.oneMinSales[last-1].StartTime()) {
				return fmt.Errorf("one minute sales times must be appended in chronological order - last time = %s, new time = %s",
					s.oneMinSales[last-1].DateTimeStr(), oms.DateTimeStr())
			}
		}

		s.oneMinSales = append(s.oneMinSales, oms)
	}

	return nil
}

func (s *Stock) String() string {
	s.Lock()
	defer s.Unlock()

	// Symbol
	str := fmt.Sprintf("Stock: %s\n", s.symbol)

	// Quote
	if len(s.quote.Symbol()) > 0 {
		str += fmt.Sprintf("Quote: %s\n", s.quote.String()[6:]) // Remove symbol, since it's shown above.
	}

	// Option Chain (expiration dates and strike price ranges)
	for i, exp := range s.optionChain {
		if i == 0 {
			str += fmt.Sprintf("Option Chain Summary:\n")
		}
		str += fmt.Sprintf("  %s\n", exp)
	}

	// Daily Prices
	length := len(s.dailyPrices)
	str += fmt.Sprintf("Daily Price Count: %d\n", length)
	if length > 0 {
		str += fmt.Sprintf("  First Day: %s\n", s.dailyPrices[0].DateStr())
		str += fmt.Sprintf("  Last Day:  %s\n", s.dailyPrices[length-1].DateStr())
		str += DailyPriceHeader()
		pricesStr := ""
		for _, dp := range s.dailyPrices {
			pricesStr = fmt.Sprintf("    %s\n", dp) + pricesStr
		}
		str += pricesStr
	}

	// One Minute Sales
	length = len(s.oneMinSales)
	str += fmt.Sprintf("One Minute Sales Count: %d\n", length)
	if length > 0 {
		str += fmt.Sprintf("    First Minute: %s\n", s.oneMinSales[0].DateTimeStr())
		str += fmt.Sprintf("    Last Minute:  %s\n", s.oneMinSales[length-1].DateTimeStr())
		str += OneMinSaleHeader()
		for _, oms := range s.oneMinSales {
			str += fmt.Sprintf("    %s\n", oms)
		}
	}

	return str
}

func SortBySymbol(stockList []*Stock) {
	sort.Slice(
		stockList,
		func(i, j int) bool {
			return stockList[i].Symbol() < stockList[j].Symbol()
		},
	)
}

func SortByLastPrice(stockList []*Stock) {
	sort.Slice(
		stockList,
		func(i, j int) bool {
			return stockList[i].LastPrice() < stockList[j].LastPrice()
		},
	)
}
