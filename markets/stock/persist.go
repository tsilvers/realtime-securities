package stock

// DailyPriceGob is the type used to persist daily price data.
type DailyPriceGob struct {
	Year   int
	Month  int
	Day    int
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume int64
}

func (dpg DailyPriceGob) ToDailyPrice() (DailyPrice, error) {
	return NewDailyPrice(dpg.Year, dpg.Month, dpg.Day, dpg.Open, dpg.Close, dpg.High, dpg.Low, dpg.Volume)
}

func (dp DailyPrice) ToGob() DailyPriceGob {
	dpg := DailyPriceGob{}
	dpg.Year = dp.date.Year()
	dpg.Month = int(dp.date.Month())
	dpg.Day = dp.date.Day()
	dpg.Open = dp.open
	dpg.Close = dp.close
	dpg.High = dp.high
	dpg.Low = dp.low
	dpg.Volume = dp.volume

	return dpg
}
