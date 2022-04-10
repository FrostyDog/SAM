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
func LaunchMarketToleranceTicker(s *kucoin.ApiService, primarySymbol string, secondarySymbol string, tradingPair string, priceMargin float64) {

	logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	log.SetOutput(logFile)

	var priceChangeList []float64
	var priceOfStart float64
	var toleranceIndicator float64 = 0.3 //0.3
	var maxChange float64
	var minChange float64
	var threshholdBuy float64  //priceOfStart - PriceMargin * priceOfStart
	var threshholdSell float64 //priceOfStart + PriceMargin * priceOfStart
	var primaryCapability int64
	var secondaryCapability int64
	var tradeAmount string = fmt.Sprintf("%v", 1)

	ticker := time.NewTicker(2 * time.Second)

	for _ = range ticker.C {
		priceInString := do.GetCurrentPrice(api.S, tradingPair)
		currentPrice, _ := strconv.ParseFloat(priceInString, 64)

		primaryCapability = calcPrimaryCapability(do.CurrencyHodlings(api.S, primarySymbol), tradeAmount)
		secondaryCapability = calcSecondaryCapability(do.CurrencyHodlings(api.S, secondarySymbol), currentPrice)

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

		if canSell && primaryCapability > 0 {
			do.MarketOrder(api.S, "sell", tradingPair, tradeAmount)
			priceChangeList = nil
			priceOfStart = 0
			log.Printf("Time to sell, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, priceOfStart, currentPrice)
			fmt.Printf("Time to sell, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, priceOfStart, currentPrice)
		}

		if canBuy && secondaryCapability > 0 {
			do.MarketOrder(api.S, "buy", tradingPair, tradeAmount)
			priceChangeList = nil
			priceOfStart = 0
			log.Printf("Time to buy, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, priceOfStart, currentPrice)
			fmt.Printf("Time to buy, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, priceOfStart, currentPrice)
		}

	}

}

func toleranceThreshhold(v float64, toleranceRate float64) float64 {
	return v - (v * toleranceRate)
}

// calculates capability based on tradeAmount (ability to sell --- floored to whole number)
func calcPrimaryCapability(avaliable float64, tradeAmount string) int64 {
	tradeAmountParsed, _ := strconv.ParseFloat(tradeAmount, 64)
	ans := avaliable / tradeAmountParsed
	return int64(math.Floor(ans))
}

// calculates capability based current Price (ability to buy --- floored to whole number)
func calcSecondaryCapability(avaliable float64, currentPrice float64) int64 {
	ans := avaliable / currentPrice
	return int64(math.Floor(ans))
}
