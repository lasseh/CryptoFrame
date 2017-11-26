package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lasseh/goi2c/devices/ledBackpack7Segment"
)

// Cryptowatch is the returning json
type Cryptowatch struct {
	Result struct {
		Price float64 `json:"price"`
	} `json:"result"`
}

var exchange = "bitfinex"

var btcURL = "https://api.cryptowat.ch/markets/" + exchange + "/btcusd/price"
var ethURL = "https://api.cryptowat.ch/markets/" + exchange + "/ethusd/price"

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
		btc := Cryptowatch{}
		getJSON(btcURL, &btc)
		// Add padding to right align the number
		btcValue := fmt.Sprintf("%4d", int(btc.Result.Price))
		display1.WriteString(btcValue)
		fmt.Println("Bitcoin Price:", btc.Result.Price)

		// Ethereum
		eth := Cryptowatch{}
		getJSON(ethURL, &eth)
		// Add padding to right align the number
		ethValue := fmt.Sprintf("%4d", int(eth.Result.Price))
		display2.WriteString(ethValue)
		fmt.Println("Etherum Price:", eth.Result.Price)

		time.Sleep(30 * time.Second)
	}

	// fmt.Println("Cleaning up...")
}

func getJSON(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}
