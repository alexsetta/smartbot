package rsi

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const period = 14

type RSI struct {
	id     string
	prices []float64
}

// NewRSI returns a new RSI struct
func NewRSI(id string) *RSI {
	return &RSI{
		id: id,
	}
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

	r.save()

	return math.Round(rsi*100) / 100
}

// get file name
func (r *RSI) fileName() string {
	return fmt.Sprintf("../files/rsi_%s.txt", strings.ToLower(r.id))
}

// save buffer to file
func (r *RSI) save() {
	buffer := ""
	for _, price := range r.prices {
		buffer += fmt.Sprintf("%f ", price)
	}

	err := os.WriteFile(r.fileName(), []byte(buffer), 0644)
	if err != nil {
		log.Println(err)
	}
}

// load buffer from file
func (r *RSI) Load() {
	buffer, err := os.ReadFile(r.fileName())
	if err != nil {
		log.Println(err)
	}
	prices := strings.Split(string(buffer), " ")
	for _, price := range prices {
		if price != "" {
			val, err := strconv.ParseFloat(price, 64)
			if err == nil {
				r.Add(val)
			}
		}
	}
}
