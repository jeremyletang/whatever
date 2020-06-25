package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}

// uncomment this to test locally

// func main() {
// 	out, err := getLeaderboard()
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return
// 	}
// 	buf, err := json.Marshal(out)
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("%s", string(buf))
// }

type Leaderboard struct {
	Traders []Trader `json:"traders"`
	Summary string   `json:"summary"`
}

type Trader struct {
	TotalUSDVal         float64 `json:"totalUsdVal"`
	USDDeployedVal      float64 `json:"usdDeployedVal"`
	PublicKey           string  `json:"publicKey"`
	BTCDeployedVal      float64 `json:"btcDeployedVal"`
	BTCVal              float64 `json:"btcVal"`
	USDDeployed         string  `json:"usdDeployed"`
	BTCDeployed         string  `json:"btcDeployed"`
	TotalUSD            string  `json:"totalUsd"`
	TotalUSDDeployedVal float64 `json:"totalUsdDeployedVal"`
	USDVal              float64 `json:"usdVal"`
	Order               int64   `json:"order"`
	TotalUSDDeployed    string  `json:"totalUsdDeployed"`
	BTC                 string  `json:"btc"`
	USD                 string  `json:"usd"`

	// added by us
	Name string `json:"name"`
}

func getLeaderboard() (*Leaderboard, error) {
	resp, err := http.Get("https://topgun-service-testnet.ops.vega.xyz/leaderboard")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	lb := Leaderboard{}
	err = json.Unmarshal([]byte(body), &lb)
	if err != nil {
		return nil, err
	}

	// first add the names
	out := []Trader{}
	for _, v := range lb.Traders {
		name, ok := keys[v.PublicKey]
		if !ok {
			v.Name = "ANONYMOUS"
		} else {
			v.Name = name
		}
		out = append(out, v)
	}
	lb.Traders = out

	// then sort
	sort.Slice(lb.Traders, func(i, j int) bool { return lb.Traders[i].Order < lb.Traders[j].Order })

	// then summary
	var i int = 1
	for _, v := range lb.Traders {
		lb.Summary = fmt.Sprintf("%v%v. %v|%v|$%v\n", lb.Summary, i, v.Name, v.PublicKey[len(v.PublicKey)-5:], v.TotalUSD)
		i++
	}

	return &lb, nil
}

func handler(ev events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST",
		"Access-Control-Allow-Headers": "*",
	}
	multiValHeaders := map[string][]string{
		"Access-Control-Allow-Methods": []string{"POST"},
		"Access-Control-Allow-Headers": []string{"*"},
	}

	if strings.EqualFold(ev.HTTPMethod, "OPTIONS") {
		// cors
		return &events.APIGatewayProxyResponse{
			StatusCode:        200,
			Headers:           headers,
			MultiValueHeaders: multiValHeaders,
		}, nil
	}

	if !strings.EqualFold(ev.HTTPMethod, "GET") {
		return nil, errors.New("method GET supported only")
	}

	// req := request{}
	// err := json.Unmarshal([]byte(ev.Body), &req)
	// if err != nil {
	// 	return nil, err
	// }

	out, err := getLeaderboard()
	if err != nil {
		return nil, err
	}
	buf, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}

	var ret string = string(buf)

	headers["Content-Type"] = "application/json"
	headers["Content-Length"] = strconv.Itoa(len(ret))

	return &events.APIGatewayProxyResponse{
		StatusCode:        200,
		Body:              ret,
		Headers:           headers,
		MultiValueHeaders: multiValHeaders,
	}, nil
}
