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

func OrderExists(s *kucoin.ApiService) bool {
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

	return len(as) != 0
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

// Median between 24h Median Price and Current price
func GetCorrelationPrice(s *kucoin.ApiService, symbol string) (correlactionPrice string) {

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
	currentPrice, _ := strconv.ParseFloat(as.Last, 64)

	var calculatedPrice float64 = (smallestPrice + highestPrice) / 2
	var resPrice float64 = (calculatedPrice + currentPrice) / 2

	return fmt.Sprintf("%v", resPrice)

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

func Get24hStats(s *kucoin.ApiService, symbol string) (stats kucoin.Stats24hrModel) {
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

	return as

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

	size := config.Size

	if sy == "" {
		sy = config.TradingPair
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

	size := config.Size

	if sy == "" {
		sy = config.TradingPair
	}

	o := kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Side: "sell", Symbol: sy, Price: price, Size: size}

	_, err := s.CreateOrder(&o)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("sell order is created")
	return "buy"
}

func MarketOrder(s *kucoin.ApiService, side string, sy string, size string) {

	if size == "" {
		size = config.Size
	}

	if sy == "" {
		sy = config.TradingPair
	}

	o := kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Type: "market", Side: side, Symbol: sy, Size: size}

	_, err := s.CreateOrder(&o)

	if err != nil {
		log.Fatal(err)
	}
}

func CurrencyHodlings(s *kucoin.ApiService, sy string) float64 {

	var resp, err = s.Accounts(sy, "")
	if err != nil {
		println("failed to fetch account info")
	}

	var info = kucoin.AccountsModel{}

	if err := resp.ReadData(&info); err != nil {
		fmt.Println("some error during reading")
	}

	v, _ := strconv.ParseFloat(info[0].Available, 64)

	return utility.RoundFloat(v, 3)
}
