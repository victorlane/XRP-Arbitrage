package main

import (
	"fmt"
	"io/ioutil"

	. "arbitrage/data"

	"gopkg.in/yaml.v2"
)

var api API

type API struct {
	Key    string `yaml:"api_key"`
	Secret string `yaml:"api_secret"`
}

func parseConfg(api *API) {
	if config, err := ioutil.ReadFile("config.yml"); err == nil {
		if err := yaml.Unmarshal(config, api); err != nil {
			panic(err)
		}
	}
}

func main() {
	parseConfg(&api)

	// bitrue, err := getAllPrices()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(bitrue.Solo)

	c, err := CreateWebsocketConnection()
	if err != nil {
		panic(err)
	}

	bookorder, err := c.GetOrdersSOLO()
	if err != nil {
		panic(err)
	}
	fmt.Println(bookorder)
}
