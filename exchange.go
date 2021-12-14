package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
)

var uri = "https://openapi.bitrue.com/"

type Client struct {
	conn     *websocket.Conn
	Response chan *Response
}

// Gets the order book for XRP / SOLO
func (c *Client) getOrdersSOLO() (bookorder *Response, err error) {
	cmd := Command{
		ID:      4,
		Command: "book_offers",
		Taker_gets: TakerGets{
			Currency: "XRP",
		},
		Taker_pays: TakerPays{
			Currency: "534F4C4F00000000000000000000000000000000", // "SOLO" currency code
			Issuer:   "rsoLo2S1kiGeCcn6hCUXVrCpGMWLrRrLZz",       // SOLO issuer
		},
		Limit: 50,
	}

	c.sendCommand(cmd.toJSON())

	sr := Response{}
	err = c.conn.ReadJSON(&sr)
	err = c.checkErr(&sr)

	if err != nil {
		return
	}

	bookorder = sr.Result.Response

	return
}

// Sends a commmand to the websocket connection
func (c *Client) sendCommand(cmd []byte) (err error) {
	if err = c.conn.WriteMessage(websocket.TextMessage, cmd); err != nil {
		return
	}
	return
}

// Checks fi there is an error in the websocket response
func (c *Client) checkErr(res *Response) (errMsg error) {
	if res.Status == "error" {
		errMsg = fmt.Errorf("[ERR:%s:%d] %s", res.Error.Error, res.Error.ErrorCode, res.ErrorMessage)
	}
	return
}

// Pings the host
func (c *Client) Ping() (err error) {

	cmd := Command{
		Command: "ping",
		ID:      1,
	}

	err = c.sendCommand(cmd.toJSON())

	if err != nil {
		return
	}

	sr := Response{}

	err = c.conn.ReadJSON(&sr)

	if err != nil {
		return
	}

	err = c.checkErr(&sr)

	if err != nil {
		return
	}

	return
}

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

// Create a websocket connection to XRPL
func createWebsocketConnection() (c Client, err error) {
	u := url.URL{
		Scheme: "ws",
		Host:   "s2.ripple.com:443",
		Path:   "/",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("go-xrp dial: ", err)
		return
	}

	c.conn = conn

	err = c.Ping()

	if err != nil {
		log.Fatal("go-xrp ping: ", err)
		return
	}

	return
}
