package testmocks

import (
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// int(0.05 * math.Pow(float64(cost), 2))
// int(math.Abs(0.05*math.Pow(float64(cost), 2)*math.Sin(float64(cost)/2)))
func BuyLicense(cost int) int {
	return int(math.Abs(0.05 * math.Pow(float64(cost), 2) * math.Sin(float64(cost)/2)))
}

func GenTreasuresList() []string {
	rand.Seed(time.Now().UnixNano())
	treasures := make([]string, 0, rand.Intn(3))
	for i := 0; i < cap(treasures); i++ {
		treasures = append(treasures, uuid.NewString())
	}

	return treasures
}

func GetPaymentForTreasure() []int {
	rand.Seed(time.Now().UnixNano())
	payment := make([]int, 0, rand.Intn(30))
	for i := 0; i < cap(payment); i++ {
		rand.Seed(time.Now().UnixNano())
		payment = append(payment, rand.Int())
	}

	return payment
}
