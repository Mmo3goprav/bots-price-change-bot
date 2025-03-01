package tradingview

type Chart struct {
	CurrentPrice          float64 `json:"current_price"`
	PriceChange           float64 `json:"price_change"`
	PriceChangePercentage float64 `json:"price_change_percentage"`
	Symbol                string  `json:"symbol"`
	Timestamp             int64   `json:"timestamp"`
	Type                  string  `json:"type"`
	Volume                float64 `json:"volume"`
}

func (c Chart) Validate() bool {
	if c.Symbol == "" {
		return false
	}

	if c.CurrentPrice == 0.0 || c.PriceChange == 0.0 {
		return false
	}

	return true
}

const (
	minSubscriptionCount = 1
)
