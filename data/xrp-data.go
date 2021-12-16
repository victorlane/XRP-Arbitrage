package data

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	Response chan *Response
}

// Gets the order book for XRP / SOLO
func (c *Client) GetOrdersSOLO() (bookorder *Response, err error) {

	// Get the order book (commands / "request" being sent to XRPL)
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

	// Send the command to XRPL via websocket
	c.sendCommand(cmd.toJSON())

	// Response struct
	sr := Response{}

	// Read the response from XRPL (from sendCommand function)
	err = c.conn.ReadJSON(&sr)

	// Check if there is an error in the response
	err = c.checkErr(&sr)

	if err != nil {
		return
	}

	// Return the response
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

// Create a websocket connection to XRPL
func CreateWebsocketConnection() (c Client, err error) {
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
