package logic

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/FrostyDog/SAM/do"
	"github.com/Kucoin/kucoin-go-sdk"
)

var coins kucoin.TickersModel

func GrowScraping(s *kucoin.ApiService) {
	logSetup()

	var targetCoin *kucoin.TickerModel
	var initialGrowth string = ""

	if targetCoin == nil {
		coins = do.GetAllCoinStats(s)
		targetCoin = iterateAndSetTargetCoin(coins)
		if targetCoin != nil {
			initialGrowth = targetCoin.ChangeRate
			log.Printf("The coins %s is bought at %s growth rate", targetCoin.Symbol, initialGrowth)
			// buy a coin
			// set a stop loss (market stop, so it will not book coins)
			// set a take profit (market stop, for the same reason above)
			// ..and idle :)
		}
	} else { //case for local testing
		currentStats := do.GetCurrentStats(s, targetCoin.Symbol)
		var sold bool = assesAndSell(currentStats, initialGrowth)

		if sold {
			targetCoin = nil
			initialGrowth = ""
		}
	}
}

func logSetup() {
	logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	log.SetOutput(logFile)
}

// Search for a target coin in all coins, returns coin and initial growth rate
func iterateAndSetTargetCoin(coins kucoin.TickersModel) *kucoin.TickerModel {

	for _, coin := range coins {
		changeRate, err := strconv.ParseFloat(coin.ChangeRate, 64)
		if err != nil {
			log.Printf("Error during converstion: %v", err)
		}

		if assessRate(changeRate) && assessForLeverage(coin.Symbol) {
			return coin
		}
	}

	return nil
}

// assesing if it is time to sell the coin
func assesAndSell(stats kucoin.Stats24hrModel, initialGrowth string) bool {
	rate, err := strconv.ParseFloat(stats.ChangeRate, 64)
	if err != nil {
		log.Printf("error when comparing coins [pasing value]: %v", err)
	}

	initialRate, err := strconv.ParseFloat(initialGrowth, 64)
	if err != nil {
		log.Printf("error when comparing coins [pasing value]: %v", err)
	}

	rateDiff := rate - initialRate

	// if rise by 10% more fix the profit
	if rateDiff > 0.1 {
		log.Printf("[PROFIT] Time to sell %s with current rate: %s", stats.Symbol, stats.ChangeRate)
		return true
	}
	// if fall by 6.5% sell to stop loss
	if rateDiff < -0.065 {
		log.Printf("[Stoploss] Time to sell %s with current rate: %s", stats.Symbol, stats.ChangeRate)
		return true
	}

	return false

}

// returns true is growRate >20%
func assessRate(rate float64) bool {
	return rate > 0.2
}

// returns false if 3L or 3S are contained
func assessForLeverage(symbol string) bool {
	symbols := strings.Split(symbol, "-")
	hasLevarage := strings.Contains(symbols[0], "3L") || strings.Contains(symbols[0], "3S")
	return !hasLevarage
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
