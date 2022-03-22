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

var targetOperation string
var targetPrice string

var transactionNotExists bool = false
var nextOperation string = "sell"

var numberOfTransaction = 0

func LaunchTicker(s *kucoin.ApiService) {

	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		currentPrice = do.GetCurrentPrice(api.S, config.DSymbol)
		transactionNotExists = do.CheckOrder(api.S)

		if transactionNotExists {
			targetPrice = do.CalculatePrice(nextOperation, currentPrice)

			if nextOperation == "sell" {
				nextOperation = do.SellCoin(api.S, "", targetPrice)
			} else {
				nextOperation = do.BuyCoin(api.S, "", targetPrice)
			}

			numberOfTransaction++
		}
		if numberOfTransaction >= 40 && nextOperation == "buy" {
			ticker.Stop()
			os.Exit(0)
		}
	}

}
