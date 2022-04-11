package main

import (
	"fmt"

	api "github.com/FrostyDog/SAM/API"
	"github.com/FrostyDog/SAM/config"
	logic "github.com/FrostyDog/SAM/logic/market_tolerance_ticker"
)

func main() {
	fmt.Println("SAM is running")

	logic.LaunchMarketToleranceTicker(api.S, config.PrimarySymbol, config.SecondarySymbol, config.PriceMargin)
}
