package logic

import (
	"fmt"
	"log"
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

// Tollerance with minimum margin on the fly model (Sell "by market")
func LaunchMarketToleranceTicker(s *kucoin.ApiService, primarySymbol string, secondarySymbol string, baseMargin float64, entryPrice float64) {

	logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	log.SetOutput(logFile)

	var tradingPair string = primarySymbol + "-" + secondarySymbol
	var priceChangeList []float64
	var startPrice float64
	var toleranceIndicator float64 = 0.25 //0.3 deafult
	var maxChange float64
	var minChange float64
	var thresholdBuy float64  //startPrice - baseMargin * startPrice
	var thresholdSell float64 //startPrice + baseMargin * startPrice
	var primaryCapability int64
	var secondaryCapability int64
	var tradeAmount string = config.TradingSize

	var dayChange float64

	ticker := time.NewTicker(2 * time.Second)

	// Execute entryPrice logic only once before the ticket --- config phase
	if entryPrice != 0 {
		startPrice = entryPrice
		thresholdSell, thresholdBuy = calcPriceThresholds(startPrice, baseMargin, dayChange)
	}

	for _ = range ticker.C {
		stats := do.GetCurrentStats(api.S, tradingPair)

		dayChange, _ = strconv.ParseFloat(stats.ChangeRate, 64)

		priceInString := stats.Last
		currentPrice, _ := strconv.ParseFloat(priceInString, 64)

		primaryHoldings, err := do.CurrencyHodlings(api.S, primarySymbol)
		if err != nil {
			log.Println("Finally Failed at primary holding")
		}
		secondaryHoldings, err := do.CurrencyHodlings(api.S, secondarySymbol)
		if err != nil {
			log.Println("Finally Failed at secondary holding")
		}

		primaryCapability = calcPrimaryCapability(primaryHoldings, tradeAmount)
		secondaryCapability = calcSecondaryCapability(secondaryHoldings, currentPrice)

		if startPrice == 0 {
			startPrice = currentPrice
			thresholdSell, thresholdBuy = calcPriceThresholds(startPrice, baseMargin, dayChange)
		}

		currentChange := utility.RoundFloat(currentPrice-startPrice, 3)

		priceChangeList = append(priceChangeList, currentChange)
		minChange, maxChange = utility.MinMax(priceChangeList)

		var canSell bool = canSell(currentChange, maxChange, toleranceIndicator, currentPrice, startPrice, thresholdSell)
		var canBuy bool = canBuy(currentChange, minChange, toleranceIndicator, currentPrice, startPrice, thresholdBuy)

		fmt.Printf("Current price is %v and currentChange is %v \n", currentPrice, currentChange)

		// For debug uncomment to track all varialbes
		// log.Printf("Current price is %v  -- Current Change: %v \n thresholdSell %v, -- thresholdBuy: %v \n Max Tollerance: %v -- MinTollerance: %v, \n\n",
		// 	currentPrice, currentChange, thresholdSell, thresholdBuy, toleranceThreshhold(maxChange, toleranceIndicator), toleranceThreshhold(minChange, toleranceIndicator))

		if canSell {
			log.Printf("Time to sell, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, startPrice, currentPrice)
			fmt.Printf("Time to sell, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, startPrice, currentPrice)
			if primaryCapability > 0 {
				do.MarketOrder(api.S, "sell", tradingPair, tradeAmount)
				log.Println("Operation done")
			}
			priceChangeList = nil
			startPrice = 0
		}

		if canBuy {
			log.Printf("Time to buy, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, startPrice, currentPrice)
			fmt.Printf("Time to buy, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, startPrice, currentPrice)
			if secondaryCapability > 0 {
				do.MarketOrder(api.S, "buy", tradingPair, tradeAmount)
				log.Println("Operation done")
			}
			priceChangeList = nil
			startPrice = 0
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
	tradingAmount, _ := strconv.ParseFloat(config.TradingSize, 64)
	ans := avaliable / currentPrice / tradingAmount
	return int64(math.Floor(ans))
}

func calcPriceThresholds(price float64, baseMargin float64, dayChange float64) (sell float64, buy float64) {

	changeIndicator := (dayChange * 100)
	changeIndicator = utility.RoundFloat(changeIndicator, 3)

	sellMargin := baseMargin
	buyMargin := baseMargin

	// every 1% is 0.001 (0.1%)
	if changeIndicator > 0 {
		sellMargin = baseMargin + changeIndicator*0.0010
	} else {
		buyMargin = baseMargin + changeIndicator*-0.0010
	}

	sell = utility.RoundFloat(price+sellMargin*price, 3)
	buy = utility.RoundFloat(price-buyMargin*price, 3)
	return sell, buy
}

func canSell(currentChange float64, maxChange float64, tolerance float64, currentPrice float64, startPrice float64, thresholdSell float64) bool {
	res := currentChange <= toleranceThreshhold(maxChange, tolerance) && currentPrice > thresholdSell
	rapidRise, _ := isRapidChange(startPrice, currentPrice)

	return res || rapidRise
}

func canBuy(currentChange float64, minChange float64, tolerance float64, currentPrice float64, startPrice float64, thresholdBuy float64) bool {
	res := currentChange >= toleranceThreshhold(minChange, tolerance) && currentPrice < thresholdBuy
	_, rapidDrop := isRapidChange(startPrice, currentPrice)

	return res || rapidDrop
}

// If change was more than 3% or -3% === do the oposite action to correct the flow to balance the capabilities
func isRapidChange(startPrice float64, currentPrice float64) (rapidRise bool, rapidDrop bool) {
	changePersantage := (currentPrice - startPrice) / startPrice * 100
	rapidRise = changePersantage >= 3
	rapidDrop = changePersantage <= -3
	return rapidRise, rapidDrop
}
