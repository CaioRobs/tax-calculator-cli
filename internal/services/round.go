package services

import "math"

func RoundTo(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}

func Round2(value float64) float64 {
	return RoundTo(value, 2)
}
