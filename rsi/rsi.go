package rsi

const period = 14

type RSI struct {
	prices []float64
}

// NewRSI returns a new RSI struct
func NewRSI() *RSI {
	return &RSI{}
}

// AppendPrice appends a new price to the prices slice
func (r *RSI) Add(price float64) {
	if len(r.prices) == (period + 1) {
		r.prices = r.prices[1:]
	}
	r.prices = append(r.prices, price)
}

// calculateRSI calculates the RSI for the given period
func (r *RSI) Calculate() float64 {
	var (
		avgGain float64
		avgLoss float64
	)

	if len(r.prices) < (period + 1) {
		return 0
	}
	start := len(r.prices) - period
	finish := len(r.prices)
	period := finish - start

	for i := start; i < finish; i++ {
		if r.prices[i] > r.prices[i-1] {
			avgGain += r.prices[i] - r.prices[i-1]
		} else {
			avgLoss += r.prices[i-1] - r.prices[i]
		}
	}

	avgGain /= float64(period)
	avgLoss /= float64(period)
	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))
	return rsi
}
