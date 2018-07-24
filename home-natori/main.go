package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type FulfillmentResponse struct {
	FulfillmentText string `json:"fulfillmentText"`
}

type FulfillmentRequest struct {
	OriginalDetectIntentRequest struct{} `json:"originalDetectIntentRequest"`
	QueryResult                 struct {
		AllRequiredParamsPresent bool     `json:"allRequiredParamsPresent"`
		DiagnosticInfo           struct{} `json:"diagnosticInfo"`
		FulfillmentMessages      []struct {
			Text struct {
				Text []string `json:"text"`
			} `json:"text"`
		} `json:"fulfillmentMessages"`
		FulfillmentText string `json:"fulfillmentText"`
		Intent          struct {
			DisplayName string `json:"displayName"`
			Name        string `json:"name"`
		} `json:"intent"`
		IntentDetectionConfidence int64  `json:"intentDetectionConfidence"`
		LanguageCode              string `json:"languageCode"`
		OutputContexts            []struct {
			LifespanCount int64  `json:"lifespanCount"`
			Name          string `json:"name"`
			Parameters    struct {
				Param string `json:"param"`
			} `json:"parameters"`
		} `json:"outputContexts"`
		Parameters struct {
			Param string `json:"param"`
		} `json:"parameters"`
		QueryText string `json:"queryText"`
	} `json:"queryResult"`
	ResponseID string `json:"responseId"`
	Session    string `json:"session"`
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

	req := FulfillmentRequest{}
	json.Unmarshal([]byte(request.Body), &req)
	fmt.Println(req.QueryResult.QueryText)
	data := FulfillmentResponse{}
	if req.QueryResult.QueryText == "GOOGLE_ASSISTANT_WELCOME" || req.QueryResult.QueryText == "おやすみ" {
		data.FulfillmentText = "<speak><audio src=\"https://s3-ap-northeast-1.amazonaws.com/usamin-star/NatoriSana/goodnight.mp3\">a</audio></speak>"
	} else {
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(84)
		data.FulfillmentText = fmt.Sprintf("<speak><audio src=\"https://s3-ap-northeast-1.amazonaws.com/usamin-star/NatoriSana/tene/ttene%03d.mp3\">a</audio></speak>", num)
	}
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
