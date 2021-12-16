package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

var uri = "https://openapi.bitrue.com"

type bitruePrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type bitruePrices struct {
	Solo float64
	XRP  float64
}

// Get price of symbol from bitrue
func getBitruePrice(symbol string) (float64, error) {
	resp, err := http.Get(uri + "/api/v1/ticker/price?symbol=" + symbol)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var ticker bitruePrice
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
func getAllBitruePrices() (bitruePrices, error) {
	var bitruePrices bitruePrices
	var err error

	bitruePrices.Solo, err = getBitruePrice("SOLOUSDT")
	if err != nil {
		return bitruePrices, err
	}

	bitruePrices.XRP, err = getBitruePrice("XRPUSDT")
	if err != nil {
		return bitruePrices, err
	}

	return bitruePrices, nil
}
