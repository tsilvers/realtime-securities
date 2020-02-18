package tradier

import (
	"encoding/json"
	"github.com/tsilvers/realtime-securities/markets"
)

const marketStatusURL = "markets/clock"

type clockJSON struct {
	Clock clockStatusJSON
}

type clockStatusJSON struct {
	State string
}

func (t *Tradier) GetMarketStatus() markets.StatusType {
	var response []byte
	var err error

	response, err = t.request(marketStatusURL)
	if err != nil {
		showRequestError(marketStatusURL, err)
		return markets.StatusError
	}

	showResponseError := showResponseErrorFunc(marketStatusURL)

	marketStatusResponse := &clockJSON{}
	if err := json.Unmarshal(response, marketStatusResponse); err != nil {
		showResponseError(err)
		return markets.StatusError
	}

	status := marketStatusResponse.Clock.State

	switch status {
	case "pre":
		return markets.StatusPre
	case "open":
		return markets.StatusOpen
	case "post":
		fallthrough
	case "closed":
		return markets.StatusClosed
	}

	return markets.StatusError
}
