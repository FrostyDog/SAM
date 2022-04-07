package logic

import (
	"os"
	"time"

	api "github.com/FrostyDog/SAM/API"
	"github.com/FrostyDog/SAM/config"
	"github.com/FrostyDog/SAM/do"

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
