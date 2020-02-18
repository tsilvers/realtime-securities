package option

import (
	"fmt"
	"sort"
	"time"
)

const MaxDaysUntilExpiration = 62

// Option holds the description of a stock option and its realtime values.
type Option struct {
	symbol    string
	expDate   time.Time
	strike    float64
	call      bool
	size      int
	open      int
	Bid       float64
	BidSize   int
	Ask       float64
	AskSize   int
	Last      float64
	Change    float64
	ChangePct float64
}

// Strike holds the call and put options for a specific strike price.
type Strike struct {
	price float64
	call  Option
	put   Option
}

func (s Strike) Price() float64 {
	return s.price
}

// Expiration holds details of all options for a specific stock and expiration date.
type Expiration struct {
	date    time.Time
	strikes []Strike
}

func NewExpiration(expDate time.Time, strikes []float64) (e Expiration, err error) {
	// Add expiration date with time set to noon.
	e.date = time.Date(expDate.Year(), expDate.Month(), expDate.Day(), 12, 0, 0, 0, time.UTC)

	// Add strike prices in ascending order.
	sort.Float64s(strikes)
	for _, strike := range strikes {
		e.strikes = append(e.strikes, Strike{price: strike})
	}

	err = e.Validate()

	return
}

func (e Expiration) Validate() error {
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 12, 0, 0, 0, time.UTC)

	// Expiration dates should not have already passed.
	// Expirations older than 2 days are errors, otherwise they are returned as a warning.
	if e.date.Before(today.Add(-48 * time.Hour)) {
		return fmt.Errorf("expiration date has passed: %s", e.date.Format("2006-01-02"))
	}
	if e.date.Before(today) {
		return fmt.Errorf("%w: %s", ExpirationWarning("expiration date has passed"), e.date.Format("2006-01-02"))
	}

	// Reject expirations with too much time remaining as a warning.
	if e.date.Sub(today) > MaxDaysUntilExpiration*24*time.Hour {
		return fmt.Errorf("%w: %s", ExpirationWarning("expiration has too much time remaining"), e.date.Format("2006-01-02"))
	}

	if len(e.strikes) == 0 {
		return fmt.Errorf("expiration date %s provided without strike prices", e.date.Format("2006-01-02"))
	}

	for i, strike := range e.strikes {
		if strike.price <= 0.0 {
			return fmt.Errorf("expiration date %s has an invalid strike price: %g", e.date, strike.price)
		}

		if i > 0 && strike.price == e.strikes[i-1].price {
			return fmt.Errorf("expiration date %s has duplicate strike prices: %g", e.date, strike.price)
		}

		if i > 0 && strike.price < e.strikes[i-1].price {
			return fmt.Errorf("strike prices are not in ascending order for expiration date %s", e.date)
		}
	}

	return nil
}

func (e Expiration) Date() time.Time {
	return e.date
}

func (e Expiration) Strikes() []Strike {
	return e.strikes
}

// ExpirationDateStr returns the expiration date as "YYYY-MM-DD".
func (e Expiration) ExpirationDateStr() string {
	return e.date.Format("2006-01-02")
}

func (e Expiration) String() (s string) {
	s = fmt.Sprintf("Date: %s", e.date.Format("2006-01-02"))
	if len(e.strikes) > 0 {
		s += fmt.Sprintf("   Strike Range: %6.1f - %6.1f", e.strikes[0].price, e.strikes[len(e.strikes)-1].price)
	}

	return
}

// ExpirationWarning indicates expiration dates which can be skipped without alerting the user.
type ExpirationWarning string

func (ew ExpirationWarning) Error() string {
	return string(ew)
}
