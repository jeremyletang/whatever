package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}

type request struct {
	Key string `json:"key"`
}

type response struct {
	Name string `json:"name"`
	Key  string `json:"key"`
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

	if !strings.EqualFold(ev.HTTPMethod, "POST") {
		return nil, errors.New("method POST supported only")
	}

	req := request{}
	err := json.Unmarshal([]byte(ev.Body), &req)
	if err != nil {
		return nil, err
	}

	if len(req.Key) <= 0 {
		return nil, errors.New("missing key field")
	}

	var res response
	for k, v := range keys {
		if strings.HasSuffix(k, req.Key) {
			res.Key = k
			res.Name = v
			break
		}
	}

	if len(res.Key) <= 0 {
		return nil, errors.New("unknown public key")
	}

	buf, err := json.Marshal(res)
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
