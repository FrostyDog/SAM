package main

import (
	"fmt"

	api "github.com/FrostyDog/SAM/API"
	"github.com/FrostyDog/SAM/logic"
)

func main() {
	fmt.Println("SAM is running")

	logic.LanchMarketToleranceTicker(api.S)
}
