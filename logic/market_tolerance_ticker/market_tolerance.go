package logic

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	api "github.com/FrostyDog/SAM/API"
	"github.com/FrostyDog/SAM/do"
	"github.com/FrostyDog/SAM/utility"

	"github.com/Kucoin/kucoin-go-sdk"
)

// Tollerance with minimum margin on the fly model (Sell "by market")
func LaunchMarketToleranceTicker(s *kucoin.ApiService, currency string, tradingPair string, priceMargin float64) {

	logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	log.SetOutput(logFile)

	var priceChangeList []float64
	var priceOfStart float64
	var toleranceIndicator float64 = 0.3 //0.3
	var maxChange float64
	var minChange float64
	var threshholdBuy float64        //priceOfStart - 0.003 * priceOfStart
	var threshholdSell float64       //priceOfStart + 0.003 * priceOfStart
	var operationIndicator int64 = 0 // -2 or +2
	var tradeAmount string = fmt.Sprintf("%v", 1)

	ticker := time.NewTicker(2 * time.Second) // usually 2 * time.Second --- 10 is for debuging
	for _ = range ticker.C {
		priceInString := do.GetCurrentPrice(api.S, tradingPair)
		currentPrice, _ := strconv.ParseFloat(priceInString, 64)

		operationIndicator = calculateCapability(do.CurrencyHodlings(api.S, currency), tradeAmount)
		if priceOfStart == 0 {
			priceOfStart = currentPrice
			threshholdSell = utility.RoundFloat(priceOfStart+priceMargin*priceOfStart, 3)
			threshholdBuy = utility.RoundFloat(priceOfStart-priceMargin*priceOfStart, 3)
		}

		currentChange := utility.RoundFloat(currentPrice-priceOfStart, 3)

		priceChangeList = append(priceChangeList, currentChange)
		minChange, maxChange = utility.MinMax(priceChangeList)

		var canSell bool = currentChange <= toleranceThreshhold(maxChange, toleranceIndicator) && currentPrice > threshholdSell
		var canBuy bool = currentChange >= toleranceThreshhold(minChange, toleranceIndicator) && currentPrice < threshholdBuy

		fmt.Printf("Current price is %v and currentChange is %v \n", currentPrice, currentChange)

		// For debug uncomment to track all varialbes
		// log.Printf("Current price is %v  -- Current Change: %v \n ThreshholdSell %v, -- ThreshholdBuy: %v \n Max Tollerance: %v -- MinTollerance: %v, \n\n",
		// 	currentPrice, currentChange, threshholdSell, threshholdBuy, toleranceThreshhold(maxChange, toleranceIndicator), toleranceThreshhold(minChange, toleranceIndicator))

		if canSell && operationIndicator > 0 {
			do.MarketOrder(api.S, "sell", tradingPair, tradeAmount)
			priceChangeList = nil
			priceOfStart = 0
			log.Printf("Time to sell, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, priceOfStart, currentPrice)
			fmt.Printf("Time to sell, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, priceOfStart, currentPrice)
		}

		if canBuy && operationIndicator < 3 {
			do.MarketOrder(api.S, "buy", tradingPair, tradeAmount)
			priceChangeList = nil
			priceOfStart = 0
			log.Printf("Time to buy, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, priceOfStart, currentPrice)
			fmt.Printf("Time to sell, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, priceOfStart, currentPrice)
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
