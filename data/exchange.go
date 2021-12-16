package data

import (
	"encoding/json"
	"net/http"
	"strconv"
)

var uri = "https://openapi.bitrue.com"

type BitruePrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type BitruePrices struct {
	Solo float64
	XRP  float64
}

// Get price of symbol from bitrue
func GetBitruePrice(symbol string) (float64, error) {
	resp, err := http.Get(uri + "/api/v1/ticker/price?symbol=" + symbol)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var ticker BitruePrice
	if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
		return 0, err
	}

	priceFloat, err := strconv.ParseFloat(ticker.Price, 64)
	if err != nil {
		return 0, err
	}

	return priceFloat, nil
}

// Get all prices from bitrue
func GetAllBitruePrices() (bitruePrices BitruePrices, err error) {
	bitruePrices.Solo, err = GetBitruePrice("SOLOUSDT")
	if err != nil {
		return
	}

	bitruePrices.XRP, err = GetBitruePrice("XRPUSDT")
	if err != nil {
		return
	}

	return
}
