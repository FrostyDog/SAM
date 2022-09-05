package logic

import (
	"fmt"
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
var initialPrice string = ""
var endTimer = make(chan bool)
var timeBombStatus bool

// update with custom short term rise calculation
var snapsCounter int = 0
var snaps models.SnapshotsContainer = models.NewSnapshotsContainter()

func GrowScraping(s *kucoin.ApiService) {
	logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	log.SetOutput(logFile)

	if targetCoin == nil {
		coins = do.GetAllCoinStats(s)
		filteredCoins := filterCoins(coins)

		snapsCounter++
		// every 15 min (10 sec * 90 = 900 sec)
		if snapsCounter == 90 {
			snapsCounter = 0
			snaps.AddSnapshotAndReplace(filteredCoins)
		}

		// if enough info - look for target
		if len(snaps) == 2 {
			targetCoin = iterateAndSetTargetCoin(snaps)
		}

		// if during 3 function above token in set
		if targetCoin != nil {
			// resete values 36h
			if !timeBombStatus {
				timeBombStatus = true
				go timeBomb(s, targetCoin)
			}
			usdCapacity := usdCapacity(s)
			do.MarketOrder(s, "buy", targetCoin.Symbol, usdCapacity, "quote")
			log.Printf("The coins %s is bought at a price of %s", targetCoin.Symbol, initialPrice)

		}
	} else { //case for local testing
		currentStats := do.GetCurrentStats(s, targetCoin.Symbol)
		// sell a coin during asses and sell
		var timeForSell bool = assesAndSell(currentStats, initialPrice)

		// clean-up before next cycle
		if timeForSell {
			coinSymbol := targetCoinSymbol(targetCoin.Symbol)
			targetCoinCapacity := targetCoinCapacity(s, coinSymbol)
			do.MarketOrder(s, "sell", targetCoin.Symbol, targetCoinCapacity, "base")
			reseteValues()
		}
	}
}

// Search for a target coin in all coins, returns coin and initial growth rate
func iterateAndSetTargetCoin(snaps models.SnapshotsContainer) *kucoin.TickerModel {

	for i, coin := range snaps[0] {
		newerData := snaps[1][i]

		priceOld, err := strconv.ParseFloat(coin.Last, 64)
		if err != nil {
			log.Printf("Error during converstion: %v", err)
		}

		priceNewer, err := strconv.ParseFloat(newerData.Last, 64)
		if err != nil {
			log.Printf("Error during converstion: %v", err)
		}

		if calcRate(priceOld, priceNewer) {
			// setting initial price to the latest (that used as the base for counting)
			initialPrice = newerData.Last
			return coin
		}
	}

	return nil
}

// returns true is growRate is big enough
func calcRate(oldPrice float64, newPrice float64) bool {

	var threshhold float64 = 1.06

	calc := newPrice / oldPrice
	// if growing rate >5% in 15 min - than target this coin
	if calc > threshhold {
		log.Printf("NewPrice was: %f and oldPrice: %f, which gives calc at %f", newPrice, oldPrice, calc)
	}
	return calc > threshhold
}

func timeBomb(s *kucoin.ApiService, targetCoin *kucoin.TickerModel) {
	select {
	case <-endTimer:
		return
	case <-time.After(36 * time.Hour):
		logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer logFile.Close()
		log.SetOutput(logFile)
		log.Printf("timer has cleared")
		sellCoin(s, targetCoin)
		reseteValues()
		return
	}

}

func sellCoin(s *kucoin.ApiService, targtetCoin *kucoin.TickerModel) {
	coinSymbol := targetCoinSymbol(targtetCoin.Symbol)
	targetCoinCapacity := targetCoinCapacity(s, coinSymbol)
	do.MarketOrder(s, "sell", targetCoin.Symbol, targetCoinCapacity, "base")
}

func usdCapacity(s *kucoin.ApiService) string {
	usdHoldings, err := do.StableCurrencyHodlings(s, "USDT")
	if err != nil {
		log.Printf("Failed at fetching USDT capacity: %s", err)
	}
	return fmt.Sprint(usdHoldings)
}

func targetCoinCapacity(s *kucoin.ApiService, targetTokenName string) string {
	targetCoinHoldings, err := do.AvaliableCurrencyHodlings(s, targetTokenName)
	if err != nil {
		log.Printf("Failed at fetching TargetCoin capacity: %s", err)
	}
	return targetCoinHoldings
}

func reseteValues() {
	targetCoin = nil
	initialPrice = ""

	timeBombStatus = false
	endTimer <- true

	snaps.ClearSnapshots()
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

// return symbol of the base currency from the target coin model
func targetCoinSymbol(symbol string) string {
	symbols := strings.Split(symbol, "-")
	return symbols[0]
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

	// if rise by 6% more fix the profit
	if priceDiff > 1.06 {
		log.Printf("[PROFIT] Time to sell %s with current price: %s", stats.Symbol, stats.Last)
		return true
	}
	// if fall by 5.5% + 5% (simulation correction) sell to stop loss
	if priceDiff < 0.895 {
		log.Printf("[STOPLOSS] Time to sell %s with current price: %s", stats.Symbol, stats.Last)
		return true
	}

	return false

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
