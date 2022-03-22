package main

import (
	api "github.com/FrostyDog/SAM/API"
	"github.com/FrostyDog/SAM/logic"
)

func main() {

	logic.LaunchTicker(api.S)
}
