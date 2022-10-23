package main

import (
	"buy-btc/bitflyer"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	getPrice        string = "価格取得"
	systemManeger   string = "System Maneger"
	buyBtcApikey    string = "buy-btc-apikey"
	buyBtcApisecret string = "buy-btc-apisecret"
	order           string = "注文"
	region          string = "ap-northeast-1"
)

// 注文方法: Limit(指値)の書い注文
// 価格->現在価格の95%
// 数量->0.001BTC
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tickerChan := make(chan *bitflyer.Ticker)
	errChan := make(chan error)
	defer close(tickerChan)
	defer close(errChan)

	log.Println(getPrice)
	go bitflyer.GetTicker(tickerChan, errChan, bitflyer.Btcjpy)
	ticker := <-tickerChan
	err := <-errChan
	if err != nil {
		return getErrorResponse(err.Error()), nil
	}

	log.Println(systemManeger)
	apiKey, err := getParameter(buyBtcApikey)
	if err != nil {
		return getErrorResponse(err.Error()), nil
	}
	apiSecret, err := getParameter(buyBtcApisecret)
	if err != nil {
		return getErrorResponse(err.Error()), nil
	}

	client := bitflyer.NewAPIClient(apiKey, apiSecret)
	price, size := bitflyer.GetBuyLogic(1, ticker, 10000.0)
	orderRes, err := bitflyer.PlaceOrderWithParams(client, price, size)
	if err != nil {
		log.Println(err)
		return getErrorResponse(err.Error()), err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("res:%+v", orderRes),
	}, nil
}

// System Manegerからパラメーターを取得する関数
func getParameter(key string) (string, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ssm.New(sess, aws.NewConfig().WithRegion(region))

	params := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}
	res, err := svc.GetParameter(params)
	if err != nil {
		return "", err
	}

	return *res.Parameter.Value, nil
}

func getErrorResponse(message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: 400,
	}
}

func main() {
	lambda.Start(handler)
}
