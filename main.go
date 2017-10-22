package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lasseh/goi2c/devices/ledBackpack7Segment"
)

// Coinbase is the response from Coinbase.com api
type Coinbase struct {
	Data struct {
		Amount   float64 `json:"amount,string"`
		Base     string  `json:"base"`
		Currency string  `json:"currency"`
	} `json:"data"`
}

var btcURL = "https://api.coinbase.com/v2/prices/BTC-USD/buy"
var ethURL = "https://api.coinbase.com/v2/prices/ETH-USD/buy"

func main() {
	fmt.Println("Initiate displays")
	// Display 1
	display1, err := ledBackpack7Segment.NewLedBackpack7Segment(1, 0x70)
	if err != nil {
		panic(err)
	}
	defer display1.Close()

	// Display 2
	display2, err := ledBackpack7Segment.NewLedBackpack7Segment(1, 0x71)
	if err != nil {
		panic(err)
	}
	defer display2.Close()

	// Init display
	display1.Begin()
	defer display1.End()
	display2.Begin()
	defer display2.End()

	// Clear the displays
	display1.Clear()
	display2.Clear()

	for {
		// Bitcoin
		btc := Coinbase{}
		getJSON(btcURL, &btc)
		// Add padding to right align the number
		btcValue := fmt.Sprintf("%4d", int(btc.Data.Amount))
		display1.WriteString(btcValue)
		fmt.Println("Bitcoin Price:", btc.Data.Amount)

		// Ethereum
		eth := Coinbase{}
		getJSON(ethURL, &eth)
		// Add padding to right align the number
		ethValue := fmt.Sprintf("%4d", int(eth.Data.Amount))
		display2.WriteString(ethValue)
		fmt.Println("Etherum Price:", eth.Data.Amount)

		time.Sleep(30 * time.Second)
	}

	// fmt.Println("Cleaning up...")
}

func getJSON(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	// Coinbase requires this header, ¯\_(ツ)_/¯
	req.Header.Add("CB-VERSION", "2017-01-01")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
