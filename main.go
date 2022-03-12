package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/google/uuid"
)

const (
	api_key    string = "6212e60d29a7d50001efccd1"
	passphrase string = "y9UU2JH9ZQUxgRjb1vHV8848DR1j17"
	secret     string = "49abfb79-4d9b-4435-9ded-ab691e734d66"
)

const dSize = "10"
const dSymbol = "XLM-USDT"

var currentPrice string

var targetOperation string
var targetPrice string

var isTransaction bool

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
		"staus":     "active",
	}

	var paginationParam = kucoin.PaginationParam{PageSize: 10, CurrentPage: 0}

	resp, err := s.Orders(params, &paginationParam)
	if err != nil {
		fmt.Println("Failed at orders")
	}

	as := kucoin.OrdersModel{}
	resp.ReadData(&as)

	log.Printf("length of the orders: %d", len(as))
}

func launchTicker(s *kucoin.ApiService) {

	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		getPrice(s, dSymbol)
	}

}

func calculateTarget() {

}

func calculatePrice(side string) {
	if side == "sell" {
		t, err := strconv.ParseFloat(currentPrice, 64)
		if err != nil {
			fmt.Println("error accured during parsing")
		}
		var p float64 = t + t*0.003
		targetPrice = fmt.Sprint(p)
	}

	if side == "buy" {
		t, err := strconv.ParseFloat(currentPrice, 64)
		if err != nil {
			fmt.Println("error accured during parsing")
		}
		var p float64 = t - t*0.003
		targetPrice = fmt.Sprint(p)
	}

}

func stopTicker(t *time.Ticker) {
	t.Stop()
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

	log.Printf("Token %s: with Price %s", as.BestAsk, as.Price)
}

func buyCoin(s *kucoin.ApiService, sy string, price string) {

	size := dSize

	o := kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Side: "buy", Symbol: sy, Price: price, Size: size}

	s.CreateOrder(&o)

}

func sellCoin(s *kucoin.ApiService, sy string, price string) {

	size := dSize

	o := kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Side: "sell", Symbol: sy, Price: price, Size: size}

	s.CreateOrder(&o)

}
