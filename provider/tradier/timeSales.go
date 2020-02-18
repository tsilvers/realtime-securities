package tradier

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tsilvers/realtime-securities/markets/stock"
	"time"
)

const timeSalesURLFmt = "markets/timesales?symbol=%s&interval=1min&start=%s&session_filter=open"

type timeSalesJSON struct {
	Series seriesJSON
}

type seriesJSON struct {
	Data []dataJSON
}

// timeSalesJSONSingle is used if only one minute is returned.
type timeSalesJSONSingle struct {
	Series seriesJSONSingle
}

type seriesJSONSingle struct {
	Data dataJSON
}

type dataJSON struct {
	Time   string
	Open   json.Number
	Close  json.Number
	High   json.Number
	Low    json.Number
	Volume json.Number
	VWAP   json.Number
}

func (t *Tradier) GetTimeSales(symbol string, start time.Time) (timeSales []stock.OneMinSale) {
	timeSales = make([]stock.OneMinSale, 0)
	var response []byte
	var err error

	url := fmt.Sprintf(timeSalesURLFmt, symbol, start.Format("2006-01-02 15:04"))
	response, err = t.request(url)
	if err != nil {
		showRequestError(url, err)
		return
	}

	showResponseError := showResponseErrorFunc(url)

	timeSalesResponse := &timeSalesJSON{}
	if err := json.Unmarshal(response, timeSalesResponse); err != nil {
		// If error on Unmarshal, check if one minute was returned.
		timeSalesResponseSingle := &timeSalesJSONSingle{}
		if err := json.Unmarshal(response, timeSalesResponseSingle); err != nil {
			showResponseError(err)
			return
		}
		// Return single quote in quotesJSON struct.
		timeSalesResponse.Series.Data = append(timeSalesResponse.Series.Data, timeSalesResponseSingle.Series.Data)
	}

	data := timeSalesResponse.Series.Data
	if len(data) == 0 {
		showResponseError(errors.New("no data in reponse"))
		return
	}

	var omsYear, omsDay, omsHour, omsMinute int
	var omsMonth time.Month
	var omsTime time.Time
	var omsOpen, omsClose, omsHigh, omsLow float64
	var omsVolume int64
	var omsVWAP float64

	for _, tsJSON := range data {
		if omsTime, err = time.Parse("2006-01-02T15:04:00", tsJSON.Time); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date/Time: %v", symbol, tsJSON.Time))
			continue
		}
		omsYear, omsMonth, omsDay = omsTime.Date()
		omsHour = omsTime.Hour()
		omsMinute = omsTime.Minute()

		if omsOpen, err = tsJSON.Open.Float64(); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date/Time: %v, Open: %v", symbol, tsJSON.Time, tsJSON.Open))
			continue
		}

		if omsClose, err = tsJSON.Close.Float64(); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date/Time: %v, Close: %v", symbol, tsJSON.Time, tsJSON.Close))
			continue
		}

		if omsHigh, err = tsJSON.High.Float64(); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date/Time: %v, High: %v", symbol, tsJSON.Time, tsJSON.High))
			continue
		}

		if omsLow, err = tsJSON.Low.Float64(); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date/Time: %v, Low: %v", symbol, tsJSON.Time, tsJSON.Low))
			continue
		}

		if omsVolume, err = tsJSON.Volume.Int64(); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date/Time: %v, Volume: %v", symbol, tsJSON.Time, tsJSON.Open))
			continue
		}

		if omsVWAP, err = tsJSON.VWAP.Float64(); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date/Time: %v, Volume-Weighted Average Price: %v", symbol, tsJSON.Time, tsJSON.VWAP))
			continue
		}

		oms, err := stock.NewOneMinSale(
			omsYear, int(omsMonth), omsDay, omsHour, omsMinute,
			omsOpen, omsClose, omsHigh, omsLow,
			omsVolume, omsVWAP,
		)
		if err != nil {
			showResponseError(err)
			continue
		}

		timeSales = append(timeSales, oms)
	}

	return
}
