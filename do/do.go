package do

import (
	"fmt"
	"log"
	"strconv"

	"github.com/FrostyDog/SAM/config"
	"github.com/FrostyDog/SAM/utility"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/google/uuid"
)

func CheckOrder(s *kucoin.ApiService) bool {
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

	return len(as) == 0
}

func CalculatePrice(side string, currentPrice string) (targetPrice string) {
	if side == "sell" {
		t, err := strconv.ParseFloat(currentPrice, 64)
		if err != nil {
			fmt.Println("error accured during parsing")
		}
		var p float64 = t + t*config.PriceMargin
		return fmt.Sprint(utility.RoundFloat(p, config.DecimalPointNumber))
	}

	if side == "buy" {
		t, err := strconv.ParseFloat(currentPrice, 64)
		if err != nil {
			fmt.Println("error accured during parsing")
		}
		var p float64 = t - t*config.PriceMargin
		return fmt.Sprint(utility.RoundFloat(p, config.DecimalPointNumber))
	}

	return

}

func GetAvarage24hPrice(s *kucoin.ApiService, symbol string) (avaragePrice string) {
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

	return fmt.Sprintf("%v", calculatedPrice)

}

func GetCurrentPrice(s *kucoin.ApiService, symbol string) (currentPrice string) {
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

	return as.Last

}

func BuyCoin(s *kucoin.ApiService, sy string, price string) (nextOperation string) {

	size := config.DSize

	if sy == "" {
		sy = config.DSymbol
	}

	o := kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Side: "buy", Symbol: sy, Price: price, Size: size}

	_, err := s.CreateOrder(&o)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("buy order is created")
	return "sell"
}

func SellCoin(s *kucoin.ApiService, sy string, price string) (nextOperation string) {

	size := config.DSize

	if sy == "" {
		sy = config.DSymbol
	}

	o := kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Side: "sell", Symbol: sy, Price: price, Size: size}

	_, err := s.CreateOrder(&o)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("sell order is created")
	return "buy"

}
