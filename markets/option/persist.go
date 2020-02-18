package option

import "time"

// ExpirationGob is the type used to persist option expiration dates and strike prices.
type ExpirationGob struct {
	Date    time.Time
	Strikes []float64
}

func (eg ExpirationGob) ToExpiration() (Expiration, error) {
	return NewExpiration(eg.Date, eg.Strikes)
}

func (e Expiration) ToGob() ExpirationGob {
	eg := ExpirationGob{}
	eg.Date = e.date
	for _, strike := range e.strikes {
		eg.Strikes = append(eg.Strikes, strike.price)
	}

	return eg
}
