package utility

import (
	"math"
)

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// MinMax in float64
func MinMax(array []float64) (float64, float64) {
	var max float64 = array[0]
	var min float64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func MinMaxSingle(lastMax float64, lastMin float64, currentChange float64) (min float64, max float64) {

	min = lastMin
	max = lastMax

	if currentChange > lastMax {
		max = currentChange
	}
	if lastMin > currentChange {
		min = currentChange
	}

	return min, max
}
