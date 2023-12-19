package utils

import "math"

func RoundFloat(number float64, precision uint8) float64 {
	precisionFactor := math.Pow10(int(precision))
	return math.Round(number*precisionFactor) / precisionFactor
}
