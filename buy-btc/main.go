package main

import (
	"buy-btc/bitflyer"
	"fmt"
	"log"
	"math"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// 注文方法: Limit(指値)の書い注文
// 価格->現在価格の95%
// 数量->0.001BTC
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	apiKey, err := getParameter("buy-btc-apikey")
	if err != nil {
		return getErrorResponse(err.Error()), nil
	}
	apisecret, err := getParameter("buy-btc-apisecret")
	if err != nil {
		return getErrorResponse(err.Error()), nil
	}

	// 現在価格取得
	ticker, err := bitflyer.GetTicker(bitflyer.Btcjpy)
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	// 購入価格取得小数点以下切り捨て
	buyPrice := RoundDecimal(ticker.Ltp * 0.95)

	order := bitflyer.Order{
		ProductCode:    bitflyer.Btcjpy.String(),
		ChildOrderType: bitflyer.Limit.String(),
		Side:           bitflyer.Buy.String(),
		Price:          buyPrice,
		Size:           0.001,
		MinuteToExpire: 4320, //3days
		TimeInForce:    bitflyer.Gtc.String(),
	}

	orderRes, err := bitflyer.PlaceOrder(&order, apiKey, apisecret)
	if err != nil {
		log.Println("------------")
		log.Println(err)
		log.Println("------------")
		return getErrorResponse(err.Error()), err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("res:%+v", orderRes),
	}, nil
}

func RoundDecimal(num float64) float64 {
	return math.Round(num)
}

// System Manegerからパラメーターを取得する関数
func getParameter(key string) (string, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ssm.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))

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
