package tradier

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tsilvers/realtime-securities/markets/quote"
	"strings"
)

const quotesURLFmt = "markets/quotes?symbols=%s"

type quotesJSON struct {
	Quotes quoteJSON
}

type quoteJSON struct {
	Quote []symbolQuoteJSON
}

// quotesJSONSingle is used if only one quote is returned.
type quotesJSONSingle struct {
	Quotes quoteJSONSingle
}

type quoteJSONSingle struct {
	Quote symbolQuoteJSON
}

type symbolQuoteJSON struct {
	Symbol    string
	TradeDate int64 `json:"trade_date"`
	PrevClose json.Number
	Change    json.Number
	ChangePct json.Number `json:"change_percentage"`
	Bid       json.Number
	BidSize   int
	Ask       json.Number
	AskSize   int
	Last      json.Number
	High      json.Number
	Low       json.Number
	Volume    int
	AvgVolume int         `json:"average_volume"`
	YearHigh  json.Number `json:"week_52_high"`
	YearLow   json.Number `json:"week_52_low"`
}

func (t *Tradier) GetQuotes(symbols []string) (quotes []quote.Quote) {
	var response []byte
	var err error

	url := fmt.Sprintf(quotesURLFmt, strings.Join(symbols, ","))
	response, err = t.request(url)
	if err != nil {
		showRequestError(url, err)
		return
	}

	showResponseError := showResponseErrorFunc(url)

	quotesResponse := &quotesJSON{}
	if err := json.Unmarshal(response, quotesResponse); err != nil {
		// If error on Unmarshal, check if one quote was returned.
		quotesResponseSingle := &quotesJSONSingle{}
		if err := json.Unmarshal(response, quotesResponseSingle); err != nil {
			showResponseError(err)
			return
		}
		// Return single quote in quotesJSON struct.
		quotesResponse.Quotes.Quote = append(quotesResponse.Quotes.Quote, quotesResponseSingle.Quotes.Quote)
	}

	quoteList := quotesResponse.Quotes.Quote
	if len(quoteList) == 0 {
		showResponseError(errors.New("no data in reponse"))
		return
	}

	var qSymbol string
	var qTradeDate int64
	var qPrevClose float64
	var qChange float64
	var qChangePct float64
	var qBid float64
	var qBidSize int
	var qAsk float64
	var qAskSize int
	var qLast float64
	var qHigh float64
	var qLow float64
	var qVolume int
	var qAvgVolume int
	var qYearHigh float64
	var qYearLow float64

	for _, qJSON := range quoteList {
		qSymbol = qJSON.Symbol

		qTradeDate = qJSON.TradeDate

		// Invalid numeric values are left as the zero value.
		// See comments on Quote.Validate()

		qPrevClose, _ = qJSON.PrevClose.Float64()

		qChange, _ = qJSON.Change.Float64()

		qChangePct, err = qJSON.ChangePct.Float64()

		qBid, _ = qJSON.Bid.Float64()

		qBidSize = qJSON.BidSize

		qAsk, _ = qJSON.Ask.Float64()

		qAskSize = qJSON.AskSize

		qLast, _ = qJSON.Last.Float64()

		qHigh, _ = qJSON.High.Float64()

		qLow, _ = qJSON.Low.Float64()

		qVolume = qJSON.Volume

		qAvgVolume = qJSON.AvgVolume

		qYearHigh, _ = qJSON.YearHigh.Float64()

		qYearLow, _ = qJSON.YearLow.Float64()

		q, err := quote.NewQuote(
			qSymbol, qTradeDate, qPrevClose, qChange, qChangePct, qBid, qBidSize, qAsk, qAskSize,
			qLast, qHigh, qLow, qVolume, qAvgVolume, qYearHigh, qYearLow)
		if err != nil {
			showResponseError(err)
			continue
		}

		quotes = append(quotes, q)
	}

	return
}
