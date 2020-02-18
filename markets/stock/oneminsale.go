package stock

import (
	"fmt"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	MarketOpenHour    = 9
	MarketOpenMinute  = 30
	MarketCloseHour   = 15
	MarketCloseMinute = 59

	maxDaysOld = 10
)

// OneMinSale holds price and volume data for one minute of stock trading activity.
type OneMinSale struct {
	startTime time.Time // Ending time of the one minute period
	open      float64
	close     float64
	high      float64
	low       float64
	volume    int64
	vwap      float64 // Volume-weighted average price
}

func NewOneMinSale(
	year, month, day, hour, min int,
	open, close, high, low float64,
	volume int64,
	vwap float64,
) (oms OneMinSale, err error) {
	oms.startTime = time.Date(year, time.Month(month), day, hour, min, 0, 0, time.UTC)
	oms.open = open
	oms.close = close
	oms.high = high
	oms.low = low
	oms.volume = volume
	oms.vwap = vwap

	err = oms.Validate()

	return
}

func (oms OneMinSale) Validate() error {
	year, month, day := time.Now().Date()
	earliestSaleTime := time.Date(year, month, day-maxDaysOld, MarketOpenHour, MarketOpenMinute, 0, 0, time.UTC)
	if oms.startTime.Before(earliestSaleTime) {
		return fmt.Errorf("one minute sales time cannot be more than %d days old", maxDaysOld)
	}

	if oms.startTime.After(time.Now()) {
		return fmt.Errorf("one minute sales time cannot be in the future")
	}

	dow := oms.startTime.Weekday()
	if dow == time.Saturday || dow == time.Sunday {
		return fmt.Errorf("one minute sales time cannot be on the weekend")
	}

	hour := oms.startTime.Hour()
	min := oms.startTime.Minute()
	if hour < MarketOpenHour || (hour == MarketOpenHour && min < MarketOpenMinute) {
		return fmt.Errorf("one minute sales time cannot be before %02d:%02d", MarketOpenHour, MarketOpenMinute)
	}
	if hour > MarketCloseHour || (hour == MarketCloseHour && min > MarketCloseMinute) {
		return fmt.Errorf("one minute sales time cannot be after %02d:%02d", MarketCloseHour, MarketCloseMinute)
	}

	if oms.open <= 0 {
		return fmt.Errorf("invalid Open price: %g", oms.open)
	}

	if oms.close <= 0 {
		return fmt.Errorf("invalid close price: %g", oms.close)
	}

	if oms.high <= 0 {
		return fmt.Errorf("invalid high price: %g", oms.high)
	}

	if oms.low <= 0 {
		return fmt.Errorf("invalid low price: %g", oms.low)
	}

	if oms.volume < 0 {
		return fmt.Errorf("invalid Volume: %d", oms.volume)
	}

	if oms.vwap < 0 {
		return fmt.Errorf("invalid vwap: %g", oms.vwap)
	}

	return nil
}

func (oms OneMinSale) StartTime() time.Time {
	return oms.startTime
}

// DateTimeStr returns the one minute sales time as "YYYY-MM-DD HH:MM".
func (oms OneMinSale) DateTimeStr() string {
	return oms.startTime.Format("2006-01-02 15:04")
}

func (oms OneMinSale) OpenClose() (float64, float64) {
	return oms.open, oms.close
}

func (oms OneMinSale) HighLow() (float64, float64) {
	return oms.high, oms.low
}

func (oms OneMinSale) Volume() int64 {
	return oms.volume
}

func (oms OneMinSale) VWAP() float64 {
	return oms.vwap
}

func (oms OneMinSale) String() string {
	p := message.NewPrinter(language.English)

	return fmt.Sprintf("%s %7.2f %7.2f %7.2f %7.2f %s %9.4f",
		oms.DateTimeStr(), oms.open, oms.close, oms.high, oms.low,
		p.Sprintf("%9d", oms.volume), oms.vwap,
	)
}

func OneMinSaleHeader() string {
	return fmt.Sprintln("                        Open   Close    High     Low    Volume      VWAP")
}
