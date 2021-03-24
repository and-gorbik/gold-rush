package testmocks

import "math"

// int(0.05 * math.Pow(float64(cost), 2))
// int(math.Abs(0.05*math.Pow(float64(cost), 2)*math.Sin(float64(cost)/2)))
func BuyLicense(cost int) int {
	return int(math.Abs(0.05 * math.Pow(float64(cost), 2) * math.Sin(float64(cost)/2)))
}
