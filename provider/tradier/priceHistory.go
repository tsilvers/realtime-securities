package tradier

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tsilvers/realtime-securities/markets/stock"
	"time"
)

const priceHistoryURLFmt = "markets/history?symbol=%s&start=%s"

type priceHistoryJSON struct {
	History historyJSON
}

type historyJSON struct {
	Day []dailyPriceJSON
}

// priceHistoryJSONSingle is used if only one day is returned.
type priceHistoryJSONSingle struct {
	History historyJSONSingle
}

type historyJSONSingle struct {
	Day dailyPriceJSON
}

type dailyPriceJSON struct {
	Date   string
	Open   json.Number
	Close  json.Number
	High   json.Number
	Low    json.Number
	Volume json.Number
}

func (t *Tradier) GetPriceHistory(symbol string, start time.Time) (dailyPrices []*stock.DailyPrice) {
	var response []byte
	var err error

	url := fmt.Sprintf(priceHistoryURLFmt, symbol, start.Format("2006-01-02"))
	response, err = t.request(url)
	if err != nil {
		showRequestError(url, err)
		return
	}

	showResponseError := showResponseErrorFunc(url)

	priceHistoryResponse := &priceHistoryJSON{}
	if err := json.Unmarshal(response, priceHistoryResponse); err != nil {
		// If error on Unmarshal, check if one day was returned.
		priceHistoryResponseSingle := &priceHistoryJSONSingle{}
		if err := json.Unmarshal(response, priceHistoryResponseSingle); err != nil {
			showResponseError(err)
			return
		}
		// Return single quote in quotesJSON struct.
		priceHistoryResponse.History.Day = append(priceHistoryResponse.History.Day, priceHistoryResponseSingle.History.Day)
	}

	day := priceHistoryResponse.History.Day
	if len(day) == 0 {
		showResponseError(errors.New("no data in reponse"))
		return
	}

	var dpYear, dpDay int
	var dpMonth time.Month
	var dpDate time.Time
	var dpOpen, dpClose, dpHigh, dpLow float64
	var dpVolume int64

	for _, dpJSON := range day {
		if dpDate, err = time.Parse("2006-01-02", dpJSON.Date); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date: %v", symbol, dpJSON.Date))
			continue
		}
		dpYear, dpMonth, dpDay = dpDate.Date()

		if dpOpen, err = dpJSON.Open.Float64(); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date: %v, Open: %v", symbol, dpJSON.Date, dpJSON.Open))
			continue
		}

		if dpClose, err = dpJSON.Close.Float64(); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date: %v, Close: %v", symbol, dpJSON.Date, dpJSON.Close))
			continue
		}

		if dpHigh, err = dpJSON.High.Float64(); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date: %v, High: %v", symbol, dpJSON.Date, dpJSON.High))
			continue
		}

		if dpLow, err = dpJSON.Low.Float64(); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date: %v, Low: %v", symbol, dpJSON.Date, dpJSON.Low))
			continue
		}

		if dpVolume, err = dpJSON.Volume.Int64(); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date: %v, Volume: %v", symbol, dpJSON.Date, dpJSON.Volume))
			continue
		}

		dp, err := stock.NewDailyPrice(dpYear, int(dpMonth), dpDay, dpOpen, dpClose, dpHigh, dpLow, dpVolume)
		if err != nil {
			showResponseError(err)
			continue
		}

		dailyPrices = append(dailyPrices, &dp)
	}

	return
}
