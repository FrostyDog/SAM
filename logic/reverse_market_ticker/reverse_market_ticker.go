package logic

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"time"

// 	api "github.com/FrostyDog/SAM/API"
// 	"github.com/FrostyDog/SAM/do"
// 	"github.com/FrostyDog/SAM/utility"

// 	"github.com/Kucoin/kucoin-go-sdk"
// )

// // Tollerance with minimum margin on the fly model (Sell "by market")
// func LaunchReverseMarketTicker(s *kucoin.ApiService, t *time.Ticker, primarySymbol string, secondarySymbol string, baseMargin float64) {

// 	logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 	defer logFile.Close()
// 	log.SetOutput(logFile)

// 	var tradingPair string = primarySymbol + "-" + secondarySymbol
// 	var startPrice float64
// 	var toleranceIndicator float64 = 0.3 //0.3 deafult
// 	var maxChange float64 = 0.00
// 	var minChange float64 = 0.00
// 	var forwardThreshold float64
// 	var backThreshold float64

// 	var primaryHoldings float64
// 	var secondaryHoldings float64

// 	for _ = range t.C {
// 		stats := do.GetCurrentStats(api.S, tradingPair)

// 		priceInString := stats.Last
// 		currentPrice, _ := strconv.ParseFloat(priceInString, 64)

// 		if startPrice == 0 {
// 			startPrice = currentPrice
// 		}

// 		if primaryHoldings == 0 && secondaryHoldings == 0 {
// 			primaryHoldings, secondaryHoldings = calcHoldings(primarySymbol, secondarySymbol)
// 		}

// 		// checking zero status
// 		if forwardThreshold == 0 && backThreshold == 0 {

// 			var secondaryCapability = calcSecondaryCapability(secondaryHoldings, currentPrice)

// 			if startPrice-startPrice*0.01 >= currentPrice {
// 				do.MarketOrder(api.S, "sell", tradingPair, fmt.Sprintf("%v", primaryHoldings))
// 				forwardThreshold = currentPrice - currentPrice*0.02
// 				backThreshold = currentPrice + currentPrice*0.005
// 			}

// 			if startPrice+startPrice*0.01 <= currentPrice {
// 				do.MarketOrder(api.S, "buy", tradingPair, secondaryCapability)
// 				forwardThreshold = currentPrice + currentPrice*0.02
// 				backThreshold = currentPrice - currentPrice*0.005
// 			}
// 		}

// 		// second stage

// 		if forwardThreshold != 0 && backThreshold != 0 {
// 			currentChange := utility.RoundFloat(currentPrice-startPrice, 3) // ex. -2 or 2
// 			minChange, maxChange = utility.MinMaxSingle(maxChange, minChange, currentChange)
// 			dropExpected := backThreshold > forwardThreshold

// 			// use switch + calculate tolerance with change () (25%)

// 			// 	switch {
// 			// 	case: dropdropExpected && currentPrice < fo
// 			// 	}

// 		}

// 		if canSell && primaryCapability > 0 {
// 			log.Printf("Time to sell, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, startPrice, currentPrice)
// 			fmt.Printf("Time to sell, current change: %v \n With Price Of start: %v and current is %v \n", currentChange, startPrice, currentPrice)
// 			do.MarketOrder(api.S, "sell", tradingPair, tradeAmount)
// 			priceChangeList = nil
// 			startPrice = 0
// 		}

// 	}

// }

// // reset(&backThreshold, &forwardThreshold, &startPrice)
// func reset(b *float64, f *float64, startPrice *float64) {
// 	*b = 0.00
// 	*f = 0.00
// 	*startPrice = 0.00
// }

// func calcHoldings(p string, s string) (float64, float64) {

// 	primaryHoldings, err := do.CurrencyHodlings(api.S, p)
// 	if err != nil {
// 		log.Println("Finally Failed at primary holding")
// 	}
// 	secondaryHoldings, err := do.CurrencyHodlings(api.S, s)
// 	if err != nil {
// 		log.Println("Finally Failed at secondary holding")
// 	}

// 	return primaryHoldings, secondaryHoldings
// }

// // calculates capability based current Price (ability to buy --- floored to whole number)
// func calcSecondaryCapability(avaliable float64, currentPrice float64) string {
// 	return fmt.Sprintf("%v", avaliable/currentPrice)
// }

// func toleranceThreshhold(v float64, toleranceRate float64) float64 {
// 	return v - (v * toleranceRate)
// }
