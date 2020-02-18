package tradier

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tsilvers/realtime-securities/markets/option"
	"time"
)

const expirationsURLFmt = "markets/options/expirations?symbol=%s&strikes=true"

type expirationsJSON struct {
	Expirations expirationJSON
}

type expirationJSON struct {
	Expiration []expirationDateJSON
}

type expirationDateJSON struct {
	Date    string
	Strikes strikesJSON
}

type strikesJSON struct {
	Strike []json.Number
}

func (t *Tradier) GetOptionExpirations(symbol string) (expirations []option.Expiration) {
	expirations = make([]option.Expiration, 0)
	var response []byte
	var err error

	url := fmt.Sprintf(expirationsURLFmt, symbol)
	response, err = t.request(url)
	if err != nil {
		showRequestError(url, err)
		return
	}

	showResponseError := showResponseErrorFunc(url)

	expirationsResponse := &expirationsJSON{}
	if err := json.Unmarshal(response, expirationsResponse); err != nil {
		showResponseError(err)
		return
	}

	data := expirationsResponse.Expirations.Expiration
	if len(data) == 0 {
		showResponseError(errors.New("no data in reponse"))
		return
	}

	var expDate time.Time
	var strikes []float64
	var expirationWarning option.ExpirationWarning

	for _, exp := range data {
		if expDate, err = time.Parse("2006-01-02", exp.Date); err != nil {
			showResponseError(err, fmt.Sprintf("Symbol: %s, Date: %s", symbol, exp.Date))
			continue
		}
		expDate = time.Date(expDate.Year(), expDate.Month(), expDate.Day(), 12, 0, 0, 0, time.UTC)

		var strikeFloat float64
		strikes = make([]float64, 0, len(exp.Strikes.Strike))
		for _, strike := range exp.Strikes.Strike {
			if strikeFloat, err = strike.Float64(); err != nil {
				showResponseError(err, fmt.Sprintf("Symbol: %s, Date: %s, Strike: %s", symbol, exp.Date, strike))
				continue
			}
			strikes = append(strikes, strikeFloat)
		}

		if expiration, err := option.NewExpiration(expDate, strikes); err != nil {
			// Don't report warnings about expiration dates.
			if !errors.As(err, &expirationWarning) {
				showResponseError(err, fmt.Sprintf("Symbol: %s, Date: %s", symbol, exp.Date))
			}
		} else {
			expirations = append(expirations, expiration)
		}
	}

	return
}
