package stock

import (
	"fmt"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Daily prices cannot be before 1/1/2010.
var earliestDate = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

// DailyPrice holds a single day's price and volume data for a security.
type DailyPrice struct {
	date   time.Time // 12pm UTC of trading day.
	open   float64
	close  float64
	high   float64
	low    float64
	volume int64
}

func NewDailyPrice(
	year, month, day int,
	open, close, high, low float64,
	volume int64,
) (dp DailyPrice, err error) {
	dp.date = time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)
	dp.open = open
	dp.close = close
	dp.high = high
	dp.low = low
	dp.volume = volume

	err = dp.Validate()

	return
}

func (dp DailyPrice) Validate() error {
	if dp.date.Before(earliestDate) {
		return fmt.Errorf("date cannot be earlier than %s", earliestDate.Format("Jan 2, 2006"))
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.UTC)
	if dp.date.After(today) {
		return fmt.Errorf("date cannot be in the future")
	}

	dow := dp.date.Weekday()
	if dow == time.Saturday || dow == time.Sunday {
		return fmt.Errorf("date cannot be on the weekend")
	}

	if dp.open <= 0 {
		return fmt.Errorf("invalid open price: %g", dp.open)
	}

	if dp.close <= 0 {
		return fmt.Errorf("invalid close price: %g", dp.close)
	}

	if dp.high <= 0 {
		return fmt.Errorf("invalid high price: %g", dp.high)
	}

	if dp.low <= 0 {
		return fmt.Errorf("invalid low price: %g", dp.low)
	}

	if dp.volume < 0 {
		return fmt.Errorf("invalid Volume: %d", dp.volume)
	}

	return nil
}

func (dp DailyPrice) Date() time.Time {
	return dp.date
}

// DateStr returns the daily price date as "YYYY-MM-DD".
func (dp DailyPrice) DateStr() string {
	return dp.date.Format("2006-01-02")
}

func (dp DailyPrice) OpenClose() (float64, float64) {
	return dp.open, dp.close
}

func (dp DailyPrice) HighLow() (float64, float64) {
	return dp.high, dp.low
}

func (dp DailyPrice) Volume() int64 {
	return dp.volume
}

func (dp DailyPrice) String() string {
	p := message.NewPrinter(language.English)

	return fmt.Sprintf("%s %7.2f %7.2f %7.2f %7.2f %s",
		dp.DateStr(), dp.open, dp.close, dp.high, dp.low,
		p.Sprintf("%12d", dp.volume),
	)
}

func DailyPriceHeader() string {
	return fmt.Sprintln("              Open   Close    High     Low       Volume")
}
