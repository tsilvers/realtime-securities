package quote

import (
	"errors"
	"fmt"
	"time"
)

const maxQuoteAgeInDays = 4

// Quote holds realtime price and volume data for a security (stock or option).
type Quote struct {
	symbol    string
	tradeDate time.Time
	prevClose float64
	change    float64
	changePct float64
	bid       float64
	bidSize   int
	ask       float64
	askSize   int
	last      float64
	high      float64
	low       float64
	volume    int
	avgVolume int
	yearHigh  float64
	yearLow   float64
}

func NewQuote(
	symbol string,
	tradeDate int64,
	prevClose float64,
	change float64,
	changePct float64,
	bid float64,
	bidSize int,
	ask float64,
	askSize int,
	last float64,
	high float64,
	low float64,
	volume int,
	avgVolume int,
	yearHigh float64,
	yearLow float64,
) (q Quote, err error) {
	q.symbol = symbol
	q.tradeDate = time.Unix(tradeDate/1000, 0)
	q.prevClose = prevClose
	q.change = change
	q.changePct = changePct
	q.bid = bid
	q.bidSize = bidSize
	q.ask = ask
	q.askSize = askSize
	q.last = last
	q.high = high
	q.low = low
	q.volume = volume
	q.avgVolume = avgVolume
	q.yearHigh = yearHigh
	q.yearLow = yearLow

	err = q.Validate()

	return
}

func (q Quote) Symbol() string {
	return q.symbol
}

func (q Quote) Last() float64 {
	return q.last
}

func (q Quote) ChangePct() float64 {
	return q.changePct
}

// Validate quote fields.
// Missing numeric values are stored as 0 and do not invalidate the quote.
// Values may be missing for several reasons, including:
//   High and Low are not provided overnight,
//   New issues do not have previous values,
//   Temporarily not given by data provider for unspecified reasons.
// Clients of quote data are expected to check for zero values.
func (q Quote) Validate() error {
	if len(q.symbol) == 0 {
		return errors.New("symbol is missing")
	}

	if q.tradeDate.IsZero() {
		return errors.New("trade date is not set")
	}

	now := time.Now()
	if q.tradeDate.After(now) {
		return fmt.Errorf("trade date cannot be in the future; trade date: %s; symbol: %s", q.tradeDate.Format("Jan 2, 2006 15:04:05"), q.symbol)
	}
	if q.tradeDate.Before(now.Add(-24 * maxQuoteAgeInDays * time.Hour)) {
		return fmt.Errorf("quote cannot be older than %d days; trade date: %s; symbol: %s", maxQuoteAgeInDays, q.tradeDate.Format("Jan 2, 2006 15:04:05"), q.symbol)
	}

	if q.prevClose < 0.0 {
		return fmt.Errorf("invalid previous close value: %g; symbol: %s", q.prevClose, q.symbol)
	}

	if q.bid < 0.0 {
		return fmt.Errorf("invalid bid value: %g; symbol: %s", q.bid, q.symbol)
	}

	if q.bidSize < 0 {
		return fmt.Errorf("invalid bid size: %d; symbol: %s", q.bidSize, q.symbol)
	}

	if q.ask < 0.0 {
		return fmt.Errorf("invalid ask value: %g; symbol: %s", q.ask, q.symbol)
	}

	if q.askSize < 0 {
		return fmt.Errorf("invalid ask size: %d; symbol: %s", q.askSize, q.symbol)
	}

	if q.last < 0.0 {
		return fmt.Errorf("invalid last value: %g; symbol: %s", q.last, q.symbol)
	}

	if q.high < 0.0 {
		return fmt.Errorf("invalid high value: %g; symbol: %s", q.high, q.symbol)
	}

	if q.low < 0.0 {
		return fmt.Errorf("invalid low value: %g; symbol: %s", q.low, q.symbol)
	}

	if q.volume < 0 {
		return fmt.Errorf("invalid volume: %d; symbol: %s", q.volume, q.symbol)
	}

	if q.avgVolume < 0 {
		return fmt.Errorf("invalid average volume: %d; symbol: %s", q.avgVolume, q.symbol)
	}

	if q.yearHigh < 0.0 {
		return fmt.Errorf("invalid 52-week high value: %g; symbol: %s", q.yearHigh, q.symbol)
	}

	if q.yearLow < 0.0 {
		return fmt.Errorf("invalid 52-week low value: %g; symbol: %s", q.yearLow, q.symbol)
	}

	return nil
}

func QuoteHeader() string {
	return fmt.Sprintln("\nStock     Last   Change   %Change\n=====     ====   ======   =======")
}

func (q Quote) String() string {
	return fmt.Sprintf(
		"%5s %8.2f %8.2f %8.2f%%",
		q.symbol, q.last, q.change, q.changePct,
	)
}
