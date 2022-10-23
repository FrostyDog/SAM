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
	"github.com/FrostyDog/SAM/utility"
	"github.com/Kucoin/kucoin-go-sdk"
)

var coins kucoin.TickersModel

var targetCoin *kucoin.TickerModel
var initialPrice string = ""
var endTimer = make(chan bool)
var timeBombStatus bool

// update with custom short term rise calculation
var snapsCounter int = 0
var snaps models.SnapshotsContainer[kucoin.TickersModel] = models.NewSnapshotsContainter[kucoin.TickersModel]()

// channels
var ticker *time.Ticker = time.NewTicker(5 * time.Second)

func GrowScraping(s *kucoin.ApiService) {
	logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	log.SetOutput(logFile)

	// put a ticker somewhere here
	if targetCoin == nil {
		snapsCounter++
		// every 15 sec (5 sec * 3 = 15 sec)
		if snapsCounter == 3 {
			coins = do.GetAllCoinStats(s)
			filteredCoins := filterCoins(coins)
			snaps.AddSnapshotAndReplace(filteredCoins)

			snapsCounter = 0
		}

		// if enough info - look for target
		if len(snaps) == 2 {
			targetCoin = iterateAndSetTargetCoin(snaps)
		}

		// if during functions above token in set
		if targetCoin != nil {
			// resete values after some time
			if !timeBombStatus {
				timeBombStatus = true
				go timeBomb(s, targetCoin)
			}
			// uncomment for real time scenario
			// targetCoinCapacity := targetCoinToBuy(s, initialPrice)
			// do.BuyCoin(s, targetCoin.Symbol, initialPrice, targetCoinCapacity)
			log.Printf("The coins %s is bought at a price of %s", targetCoin.Symbol, initialPrice)

		}
	} else { //case for local testing
		currentPrice := 
		// sell a coin during asses and sell
		var sold bool = assesAndSell(s, currentPrice, initialPrice)

		// clean-up before next cycle
		if sold {
			reseteValues()
			endTimeBomb()
		}
	}
}

// Search for a target coin in all coins, returns coin and initial growth rate
func iterateAndSetTargetCoin(snaps models.SnapshotsContainer[kucoin.TickersModel]) *kucoin.TickerModel {

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

	var threshhold float64 = 1.02  //grow rate
	var falseRateCap float64 = 1.4 // eliminating every result that is bigger than 3x (high chance it is false number)

	calc := newPrice / oldPrice
	okToBuy := calc > threshhold && calc < falseRateCap
	// if growing rate is more than threshhold - than target this coin
	if okToBuy {
		log.Printf("NewPrice was: %f and oldPrice: %f, which gives calc at %f", newPrice, oldPrice, calc)
	}
	return okToBuy
}

func timeBomb(s *kucoin.ApiService, targetCoin *kucoin.TickerModel) {
	select {
	case <-endTimer:
		return
	case <-time.After(20 * time.Hour):
		logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer logFile.Close()
		log.SetOutput(logFile)
		log.Printf("timer has cleared")
		// uncomment for real time scenario
		// sellCoinMarket(s, targetCoin)
		reseteValues()
		return
	}

}

func sellCoinMarket(s *kucoin.ApiService, targtetCoin *kucoin.TickerModel) {
	targetCoinCapacity := targetCoinCapacity(s, targtetCoin.Symbol)
	do.MarketOrder(s, "sell", targetCoin.Symbol, targetCoinCapacity, "base")
}

// return a float64 of avaliable USDT in account
func usdCapacity(s *kucoin.ApiService) float64 {
	usdHoldings, err := do.CurrencyHodlings(s, "USDT")
	if err != nil {
		log.Printf("Failed at fetching USDT capacity: %s", err)
	}
	return usdHoldings
}

// calculates an amount of base currency that could be bought with usdt capasity.
func targetCoinToBuy(s *kucoin.ApiService, price string) string {
	usdCapacity := usdCapacity(s)
	amountToBuy := usdCapacity / utility.StringToFloat64(initialPrice)

	amountToBuy = utility.RoundFloat(amountToBuy, 3)

	return fmt.Sprint(amountToBuy)
}

func targetCoinCapacity(s *kucoin.ApiService, symbol string) string {
	coinSymbol := targetCoinSymbol(symbol)
	targetCoinHoldings, err := do.CurrencyHodlings(s, coinSymbol)
	if err != nil {
		log.Printf("Failed at fetching TargetCoin capacity: %s", err)
	}
	return fmt.Sprint(targetCoinHoldings)
}

func reseteValues() {
	targetCoin = nil
	initialPrice = ""

	snaps.ClearSnapshots()
}

func endTimeBomb() {
	timeBombStatus = false
	endTimer <- true
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

// assesing and selling the coin
func assesAndSell(s *kucoin.ApiService, currentPrice string, initialPrice string) bool {
	price, err := strconv.ParseFloat(currentPrice, 64)
	if err != nil {
		log.Printf("error when parsing current price: %v", err)
	}

	initPrice, err := strconv.ParseFloat(initialPrice, 64)
	if err != nil {
		log.Printf("error when parsing initial price: %v", err)
	}

	priceDiff := price / initPrice

	// if rise by 3.5% more fix the profit
	if priceDiff > 1.035 {
		// uncomment for real time scenario
		// targetCoinCapacity := targetCoinCapacity(s, stats.Symbol)
		// do.SellCoin(s, stats.Symbol, stats.Last, targetCoinCapacity)
		log.Printf("[PROFIT] Time to sell %s coin with current price: %s", targetCoin.Symbol, currentPrice)
		return true
	}
	// if fall by 5.5% (simulation correction) sell to stop loss
	if priceDiff < 0.945 {
		// uncomment for real time scenario
		// sellCoinMarket(s, targetCoin)
		log.Printf("[STOPLOSS] Time to sell %s coin with aprx price: %s", targetCoin.Symbol, currentPrice)
		return true
	}

	return false

}
