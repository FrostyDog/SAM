package main

import (
	"fmt"

	api "github.com/FrostyDog/SAM/API"
	logic "github.com/FrostyDog/SAM/logic/market_tolerance"
)

func main() {
	fmt.Println("SAM is running")

	logic.LaunchMarketToleranceTicker(api.S)
}
