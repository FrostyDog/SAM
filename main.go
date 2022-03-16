package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/FrostyDog/SAM/utility"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/google/uuid"
)

const (
	api_key    string = "6212e60d29a7d50001efccd1"
	passphrase string = "y9UU2JH9ZQUxgRjb1vHV8848DR1j17"
	secret     string = "49abfb79-4d9b-4435-9ded-ab691e734d66"
)

const dSize = "100"
const dSymbol = "XLM-USDT"

var currentPrice string

var targetOperation string
var targetPrice string

var transactionNotExists bool = false
var nextOperation string = "sell"

var numberOfTransaction = 0

func main() {

	s := kucoin.NewApiService(
		kucoin.ApiBaseURIOption("https://api.kucoin.com"),
		kucoin.ApiKeyOption(api_key),
		kucoin.ApiSecretOption(secret),
		kucoin.ApiPassPhraseOption(passphrase),
		kucoin.ApiKeyVersionOption(kucoin.ApiKeyVersionV2))

	launchTicker(s)

}

func checkOrder(s *kucoin.ApiService) {
	var params = map[string]string{
		"tradeType": "TRADE",
		"status":    "active",
	}

	var paginationParam = kucoin.PaginationParam{PageSize: 10, CurrentPage: 1}

	resp, err := s.Orders(params, &paginationParam)
	if err != nil {
		fmt.Println("Failed at orders")
	}

	as := kucoin.OrdersModel{}
	_, err = resp.ReadPaginationData(&as) // put variable instead of blank to see pagination/page resutlss
	if err != nil {
		fmt.Println("Failed at reading pagination")
	}

	transactionNotExists = len(as) == 0
}

func launchTicker(s *kucoin.ApiService) {

	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		getPrice(s, dSymbol)
		checkOrder(s)

		if transactionNotExists {
			calculatePrice(nextOperation)

			if nextOperation == "sell" {
				sellCoin(s, "", targetPrice)
			} else {
				buyCoin(s, "", targetPrice)
			}

			numberOfTransaction++
		}

		if numberOfTransaction == 10 {
			ticker.Stop()
			os.Exit(0)
		}
	}

}

func calculatePrice(side string) {
	if side == "sell" {
		t, err := strconv.ParseFloat(currentPrice, 64)
		if err != nil {
			fmt.Println("error accured during parsing")
		}
		var p float64 = t + t*0.003
		targetPrice = fmt.Sprint(utility.RoundFloat(p, 5))
	}

	if side == "buy" {
		t, err := strconv.ParseFloat(currentPrice, 64)
		if err != nil {
			fmt.Println("error accured during parsing")
		}
		var p float64 = t - t*0.0025
		targetPrice = fmt.Sprint(utility.RoundFloat(p, 5))
	}

}

func getPrice(s *kucoin.ApiService, symbol string) {
	rsp, err := s.TickerLevel1(symbol)
	if err != nil {
		fmt.Println("error in account")
		return
	}

	as := kucoin.TickerLevel1Model{}
	if err := rsp.ReadData(&as); err != nil {
		fmt.Println("some error during reading")
		return
	}

	currentPrice = as.Price
}

func buyCoin(s *kucoin.ApiService, sy string, price string) {

	size := dSize

	if sy == "" {
		sy = dSymbol
	}

	o := kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Side: "buy", Symbol: sy, Price: price, Size: size}

	_, err := s.CreateOrder(&o)

	if err != nil {
		log.Fatal(err)
	}

	if err == nil {
		fmt.Println("buy operation done")
		nextOperation = "sell"
	}

}

func sellCoin(s *kucoin.ApiService, sy string, price string) {

	size := dSize

	if sy == "" {
		sy = dSymbol
	}

	o := kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Side: "sell", Symbol: sy, Price: price, Size: size}

	_, err := s.CreateOrder(&o)

	if err != nil {
		log.Fatal(err)
	}

	if err == nil {
		fmt.Println("sell operation done")
		nextOperation = "buy"
	}

}
