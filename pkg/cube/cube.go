package cube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
)

// load markets and asset types from cube
// https://staging.cube.exchange/ir/v0/markets

var ProdMarketsURL string = "https://api.cube.exchange/ir/v0/markets"
var StagingMarketsURL string = "https://staging.cube.exchange/ir/v0/markets"

var cachedMarketsFile string = "data/markets.json"

// Market represents a single market object.
type Market struct {
	MarketID        int64   `json:"marketId"`
	TransactTime    int64   `json:"transactTime"`
	BidPrice        float64 `json:"bidPrice"`
	BidQuantity     int     `json:"bidQuantity"`
	AskPrice        float64 `json:"askPrice"`
	AskQuantity     int     `json:"askQuantity"`
	LastPrice       float64 `json:"lastPrice"`
	Rolling24hPrice float64 `json:"rolling24hPrice"`
}

func ReadMarketsFromFile(environment string) ([]Market, error) {
	// Read the JSON data from the file
	data, err := os.ReadFile(environment + "-" + cachedMarketsFile)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a slice of Market structs
	var markets []Market
	if err := json.Unmarshal(data, &markets); err != nil {
		return nil, err
	}

	return markets, nil
}

func GetMarketsFromAPI(environment string) ([]Market, error) {
	var url string
	if environment == "production" {
		url = ProdMarketsURL
	} else if environment == "staging" {
		url = StagingMarketsURL
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request failed with status code %d", resp.StatusCode)
	}

	// Decode the JSON response into a slice of Market structs
	var markets []Market
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&markets); err != nil {
		return nil, err
	}

	return markets, nil
}

func CheckMarketsDifference(environment string) (bool, error) {
	// Read markets from the file
	fileMarkets, err := ReadMarketsFromFile(environment)
	if err != nil {
		return false, err
	}

	// Get markets from the API
	apiMarkets, err := GetMarketsFromAPI(environment)
	if err != nil {
		return false, err
	}

	// Compare the two slices of markets
	return !reflect.DeepEqual(fileMarkets, apiMarkets), nil
}
