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
		var p float64 = t + t*config.BaseMargin
		return fmt.Sprint(utility.RoundFloat(p, config.DecimalPointNumber))
	}

	if side == "buy" {
		t, err := strconv.ParseFloat(currentPrice, 64)
		if err != nil {
			fmt.Println("error accured during parsing")
		}
		var p float64 = t - t*config.BaseMargin
		return fmt.Sprint(utility.RoundFloat(p, config.DecimalPointNumber))
	}

	return

}

// Median between 24h Median Price and Current price
func GetCorrelationPrice(s *kucoin.ApiService, symbol string) (correlactionPrice string) {

	rsp, err := s.Stats24hr(symbol)
	if err != nil {
		log.Printf("error in account: %s", err)
		return
	}

	as := kucoin.Stats24hrModel{}
	if err := rsp.ReadData(&as); err != nil {
		log.Printf("some error during reading: %v", err)
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

	var rsp *kucoin.ApiResponse
	var err error

	for {
		rsp, err = s.Stats24hr(symbol)
		if err == nil {
			break
		} else {
			log.Printf("[Retrying] Error in getting Stats %v", err)
		}
	}

	stats = kucoin.Stats24hrModel{}
	if err := rsp.ReadData(&stats); err != nil {
		fmt.Println("some error during reading")
		return
	}

	return stats

}

func BuyCoin(s *kucoin.ApiService, sy string, price string) (nextOperation string) {

	size := config.TradingSize

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

	size := config.TradingSize

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

func MarketOrder(s *kucoin.ApiService, side string, sy string, size string, baseOrQuote string) {

	if baseOrQuote == "" {
		baseOrQuote = "base"
	}

	if size == "" {
		size = config.TradingSize
	}

	if sy == "" {
		sy = config.TradingPair
	}

	o := kucoin.CreateOrderModel{}

	// if set to "base" sell size with base currency else with quote(second placed)
	if baseOrQuote == "base" {
		o = kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Type: "market", Side: side, Symbol: sy, Size: size}
	} else {
		o = kucoin.CreateOrderModel{ClientOid: uuid.New().String(), Type: "market", Side: side, Symbol: sy, Funds: size}
	}
	res, err := s.CreateOrder(&o)

	if res.Code != "200000" {
		log.Printf("error is market order. response: %v", res)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func CurrencyHodlings(s *kucoin.ApiService, sy string) (holdings float64, err error) {

	var resp *kucoin.ApiResponse

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Recovered] %v", r)
			holdings, err = CurrencyHodlings(s, sy)
		}
	}()

	for {
		resp, err = s.Accounts(sy, "")
		if err != nil {
			log.Printf("[Retrying] Error in accounts %v", err)
		} else if resp.Code != "200000" {
			log.Printf("[Retrying] KuCoin internal error in accounts %v", err)
		} else {
			break
		}
	}

	var info = kucoin.AccountsModel{}

	if err := resp.ReadData(&info); err != nil {
		log.Printf("Error in reading accounts %v", err)
	}

	v, err := strconv.ParseFloat(info[0].Available, 64)

	// reserving 0.66% of all amount for transactional fees
	v = v - (v / 166)

	// flooring to 0.003 number (should work for most of the orders)
	holdings = utility.RoundFloat(v, 3)

	return holdings, err
}

func GetAllCoinStats(s *kucoin.ApiService) kucoin.TickersModel {

	var rsp *kucoin.ApiResponse
	var respErr error
	var allCoinsData kucoin.TickersResponseModel

	for {
		rsp, respErr = s.Tickers()
		if respErr != nil {
			log.Printf("[Retrying] Error in tickers %v", respErr)
		} else if rsp.Code != "200000" {
			log.Printf("[Retrying] KuCoin internal error in tickers %v", respErr)
		} else {
			break
		}
	}

	allCoinsData = kucoin.TickersResponseModel{}
	if err := rsp.ReadData(&allCoinsData); err != nil {
		fmt.Println("Error during reading all coins tickers")
	}

	return allCoinsData.Tickers
}
