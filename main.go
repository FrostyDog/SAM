package main

import (
	"fmt"
	"log"

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

func main() {

	s := kucoin.NewApiService(
		kucoin.ApiBaseURIOption("https://api.kucoin.com"),
		kucoin.ApiKeyOption(api_key),
		kucoin.ApiSecretOption(secret),
		kucoin.ApiPassPhraseOption(passphrase),
		kucoin.ApiKeyVersionOption(kucoin.ApiKeyVersionV2))

	// ticker := time.NewTicker(3 * time.Second)
	// for _ = range ticker.C {
	// 	getPrice(s, dSymbol)
	// }

	buyCoin(s, dSymbol, "0.1800")

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

	log.Printf("Token %s: with Price %s", as.BestAsk, as.Price)
}

func buyCoin(s *kucoin.ApiService, sy string, price string) {

	size := dSize

	o := kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Side: "buy", Symbol: sy, Price: price, Size: size}

	s.CreateOrder(&o)

}

func sellCoin() {

}
