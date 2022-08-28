package logic

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FrostyDog/SAM/do"
	"github.com/FrostyDog/SAM/models"
	"github.com/Kucoin/kucoin-go-sdk"
)

var coins kucoin.TickersModel

var targetCoin *kucoin.TickerModel
var initialGrowth string = ""
var initialPrice string = ""
var endTimer = make(chan bool)

// update with custom short term rise calculation

var container models.SnapshotsContainer

func GrowScraping(s *kucoin.ApiService) {
	logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	log.SetOutput(logFile)

	if targetCoin == nil {
		coins = do.GetAllCoinStats(s)
		filteredCoins := filterCoins(coins)
		targetCoin = iterateAndSetTargetCoin(filteredCoins)
		if targetCoin != nil {
			// resete values 36h
			go timeBomb(targetCoin)
			initialGrowth = targetCoin.ChangeRate
			initialPrice = targetCoin.Last
			log.Printf("The coins %s is bought at %s growth rate with a price of %s", targetCoin.Symbol, initialGrowth, initialPrice)
			// buy a coin
			// set a stop loss (market stop, so it will not book coins)
			// set a take profit (market stop, for the same reason above)
			// ..and idle :)
		}
	} else { //case for local testing
		currentStats := do.GetCurrentStats(s, targetCoin.Symbol)
		var sold bool = assesAndSell(currentStats, initialPrice)

		// clean-up before next cycle
		if sold {
			endTimer <- true
			reseteValues()
		}
	}
}

// Search for a target coin in all coins, returns coin and initial growth rate
func iterateAndSetTargetCoin(filteredCoins kucoin.TickersModel) *kucoin.TickerModel {

	for _, coin := range filteredCoins {
		changeRate, err := strconv.ParseFloat(coin.ChangeRate, 64)
		if err != nil {
			log.Printf("Error during converstion: %v", err)
		}

		if assessRate(changeRate) {
			return coin
		}
	}

	return nil
}

func timeBomb(coins *kucoin.TickerModel) {
	select {
	case <-endTimer:
		return
	case <-time.After(36 * time.Hour):
		logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer logFile.Close()
		log.SetOutput(logFile)
		log.Printf("timer has cleared")
		// TODO: market sell here at current price
		reseteValues()
		return
	}

}

func reseteValues() {
	targetCoin = nil
	initialGrowth = ""
	initialPrice = ""
}

// filter coin pair to the USDT pairs only + without levarage
func filterCoins(coins kucoin.TickersModel) kucoin.TickersModel {
	var filteredCoins kucoin.TickersModel

	for _, coin := range coins {
		symbols := strings.Split(coin.Symbol, "-")
		hasLevarage := strings.Contains(symbols[0], "3L") || strings.Contains(symbols[0], "3S")
		if symbols[1] == "USDT" && !hasLevarage {
			filteredCoins = append(filteredCoins, coin)
		}
	}

	return filteredCoins
}

// assesing if it is time to sell the coin
func assesAndSell(stats kucoin.Stats24hrModel, initialPrice string) bool {
	price, err := strconv.ParseFloat(stats.Last, 64)
	if err != nil {
		log.Printf("error when parsing current price: %v", err)
	}

	initPrice, err := strconv.ParseFloat(initialPrice, 64)
	if err != nil {
		log.Printf("error when parsing initial price: %v", err)
	}

	priceDiff := price / initPrice

	// if rise by 10% more fix the profit
	if priceDiff > 1.1 {
		log.Printf("[PROFIT] Time to sell %s with current price: %s", stats.Symbol, stats.Last)
		return true
	}
	// if fall by 6.5% sell to stop loss
	if priceDiff < 0.945 {
		log.Printf("[STOPLOSS] Time to sell %s with current price: %s", stats.Symbol, stats.Last)
		return true
	}

	return false

}

// returns true is growRate >20%
func assessRate(rate float64) bool {
	return rate > 0.2
}

// compare growsRate and returns coin with largest growth Rate
// func compareCoins(previousCoin *kucoin.TickerModel, laterCoin *kucoin.TickerModel) *kucoin.TickerModel {

// 	changeRate1, err := strconv.ParseFloat(previousCoin.ChangeRate, 64)
// 	if err != nil {
// 		log.Printf("error when comparing coins [pasing value]: %v", err)
// 	}

// 	changeRate2, err := strconv.ParseFloat(laterCoin.ChangeRate, 64)
// 	if err != nil {
// 		log.Printf("error when comparing coins [pasing value]: %v", err)
// 	}

// 	if changeRate1-changeRate2 > 0 {
// 		return previousCoin
// 	} else {
// 		return laterCoin
// 	}
// }
