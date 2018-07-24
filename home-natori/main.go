package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type FulfillmentResponse struct {
	FulfillmentText string `json:"fulfillmentText"`
}

func JSONSafeMarshal(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	return b, err
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received body: ", request.Body)

  data := FulfillmentResponse{"<speak><audio src=\"https://s3-ap-northeast-1.amazonaws.com/usamin-star/NatoriSana/goodnight.mp3\">a</audio></speak>"}

	res, err := JSONSafeMarshal(data)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return events.APIGatewayProxyResponse{Body: "JSON Marshal error:", StatusCode: 422}, nil
	}

	return events.APIGatewayProxyResponse{Body: html.UnescapeString(string(res)), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
