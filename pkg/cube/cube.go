package cube

import (
	"bufio"
	"fmt"
	"net/http"
)

// load markets and asset types from cube
// https://staging.cube.exchange/ir/v0/markets

var prodMarketsURL string = "https://api.cube.exchange/ir/v0/markets"
var stagingMarketsURL string = "https://staging.cube.exchange/ir/v0/markets"

func LoadMarkets(environment string) {
	var url string
	if environment == "production" {
		url = prodMarketsURL
	} else if environment == "staging" {
		url = stagingMarketsURL
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}