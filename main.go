package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/FrostyDog/SAM/config"
	"github.com/FrostyDog/SAM/utility"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/google/uuid"
)

var avaragePrice string
var currentPrice string

var targetOperation string
var targetPrice string

var transactionNotExists bool = false
var nextOperation string = "sell"

var numberOfTransaction = 0

func main() {

	s := kucoin.NewApiService(
		kucoin.ApiBaseURIOption("https://api.kucoin.com"),
		kucoin.ApiKeyOption(config.Api_key),
		kucoin.ApiSecretOption(config.Secret),
		kucoin.ApiPassPhraseOption(config.Passphrase),
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
		log.Fatal(err)
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
		getCurrentPrice(s, config.DSymbol)
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
		var p float64 = t + t*config.PriceMargin
		targetPrice = fmt.Sprint(utility.RoundFloat(p, config.DecimalPointNumber))
	}

	if side == "buy" {
		t, err := strconv.ParseFloat(currentPrice, 64)
		if err != nil {
			fmt.Println("error accured during parsing")
		}
		var p float64 = t - t*config.PriceMargin
		targetPrice = fmt.Sprint(utility.RoundFloat(p, config.DecimalPointNumber))
	}

}

func getAvarage24hPrice(s *kucoin.ApiService, symbol string) {
	rsp, err := s.Stats24hr(symbol)
	if err != nil {
		fmt.Println("error in account")
		return
	}

	as := kucoin.Stats24hrModel{}
	if err := rsp.ReadData(&as); err != nil {
		fmt.Println("some error during reading")
		return
	}

	highestPrice, _ := strconv.ParseFloat(as.High, 64)
	smallestPrice, _ := strconv.ParseFloat(as.Low, 64)

	var calculatedPrice float64 = (smallestPrice + highestPrice) / 2

	fmt.Printf("This is calculated price %+v \n", as)
	fmt.Printf("This is calculated price %v", calculatedPrice)

	avaragePrice = fmt.Sprintf("%v", calculatedPrice)

}

func getCurrentPrice(s *kucoin.ApiService, symbol string) {
	rsp, err := s.Stats24hr(symbol)
	if err != nil {
		fmt.Println("error in account")
		return
	}

	as := kucoin.Stats24hrModel{}
	if err := rsp.ReadData(&as); err != nil {
		fmt.Println("some error during reading")
		return
	}

	currentPrice = as.Last

}

func buyCoin(s *kucoin.ApiService, sy string, price string) {

	size := config.DSize

	if sy == "" {
		sy = config.DSymbol
	}

	o := kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Side: "buy", Symbol: sy, Price: price, Size: size}

	_, err := s.CreateOrder(&o)

	if err != nil {
		log.Fatal(err)
	}

	if err == nil {
		fmt.Println("buy order is created")
		nextOperation = "sell"
	}

}

func sellCoin(s *kucoin.ApiService, sy string, price string) {

	size := config.DSize

	if sy == "" {
		sy = config.DSymbol
	}

	o := kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Side: "sell", Symbol: sy, Price: price, Size: size}

	_, err := s.CreateOrder(&o)

	if err != nil {
		log.Fatal(err)
	}

	if err == nil {
		fmt.Println("sell order is created")
		nextOperation = "buy"
	}

}
