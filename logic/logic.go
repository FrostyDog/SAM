package logic

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	api "github.com/FrostyDog/SAM/API"
	"github.com/FrostyDog/SAM/config"
	"github.com/FrostyDog/SAM/do"
	"github.com/FrostyDog/SAM/utility"

	"github.com/Kucoin/kucoin-go-sdk"
)

var avaragePrice string
var currentPrice string

var targetPrice string

var transactionNotExists bool = false
var nextOperation string = "sell"

var numberOfTransaction = 0

// Takes current price as the base, and form selling/buying from it.
func LaunchBasicTicker(s *kucoin.ApiService) {

	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		currentPrice = do.GetCurrentPrice(api.S, config.DSymbol)
		transactionNotExists = !do.OrderExists(api.S)

		if transactionNotExists {
			targetPrice = do.CalculatePrice(nextOperation, currentPrice)

			if nextOperation == "sell" {
				nextOperation = do.SellCoin(api.S, "", targetPrice)
			} else {
				nextOperation = do.BuyCoin(api.S, "", targetPrice)
			}

			numberOfTransaction++
		}
		if numberOfTransaction >= 20 && nextOperation == "buy" {
			ticker.Stop()
			os.Exit(0)
		}
	}

}

// Takes correlactionPrice as a base
func LaounchCorrelationTicker(s *kucoin.ApiService) {

	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		currentPrice = do.GetCorrelationPrice(api.S, config.DSymbol)
		transactionNotExists = !do.OrderExists(api.S)

		if transactionNotExists {
			targetPrice = do.CalculatePrice(nextOperation, currentPrice)

			if nextOperation == "sell" {
				nextOperation = do.SellCoin(api.S, "", targetPrice)
			} else {
				nextOperation = do.BuyCoin(api.S, "", targetPrice)
			}

			numberOfTransaction++
		}
		if numberOfTransaction >= 20 && nextOperation == "buy" {
			ticker.Stop()
			os.Exit(0)
		}
	}

}

// Tollerance with minimum margin on the fly model (Sell "by market")
func LanchMarketToleranceTicker(s *kucoin.ApiService) {

	var priceChangeList []float64
	var priceOfStart float64
	var toleranceIndicator float64 = 0.3 //0.3
	var maxChange float64
	var minChange float64
	var threshholdBuy float64        //priceOfStart - 0.003 * priceOfStart
	var threshholdSell float64       //priceOfStart + 0.003 * priceOfStart
	var operationIndicator int64 = 0 // -2 or +2
	var tradeAmount string = fmt.Sprintf("%v", 0.5)

	ticker := time.NewTicker(2 * time.Second)
	for _ = range ticker.C {
		priceInString := do.GetCurrentPrice(api.S, config.DSymbol)
		currentPrice, _ := strconv.ParseFloat(priceInString, 64)

		operationIndicator = calculateCapability(do.CurrencyHodlings(api.S, "SOL"), tradeAmount)
		if priceOfStart == 0 {
			priceOfStart = currentPrice
			threshholdSell = utility.RoundFloat(priceOfStart+0.0035*priceOfStart, 3)
			threshholdBuy = utility.RoundFloat(priceOfStart-0.0035*priceOfStart, 3)
		}

		currentChange := utility.RoundFloat(currentPrice-priceOfStart, 3)

		priceChangeList = append(priceChangeList, currentChange)
		minChange, maxChange = utility.MinMax(priceChangeList)

		var canSell bool = currentChange <= toleranceThreshhold(maxChange, toleranceIndicator) && currentPrice > threshholdSell
		var canBuy bool = currentChange >= toleranceThreshhold(minChange, toleranceIndicator) && currentPrice < threshholdBuy

		fmt.Printf("Current price is %v and currentChange is %v \n", currentPrice, currentChange)

		if canSell && operationIndicator > 0 {
			do.MarketOrder(api.S, "sell", config.DSymbol, tradeAmount)
			priceChangeList = nil
			priceOfStart = 0
			fmt.Printf("Time to sell, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, priceOfStart, currentPrice)
		}

		if canBuy && operationIndicator < 2 {
			do.MarketOrder(api.S, "buy", config.DSymbol, tradeAmount)
			priceChangeList = nil
			priceOfStart = 0
			fmt.Printf("Time to buy, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, priceOfStart, currentPrice)
		}

	}

}

func toleranceThreshhold(v float64, toleranceRate float64) float64 {
	return v - (v * toleranceRate)
}

func calculateCapability(avaliable float64, tradeAmount string) int64 {
	tradeAmountParsed, _ := strconv.ParseFloat(tradeAmount, 64)
	ans := avaliable / tradeAmountParsed
	return int64(math.Round(ans))
}
