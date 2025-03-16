package main

import (
	"context"
	"log"

	"gold-rush/internal/server"
)

func main() {
	handlers := server.Init("http://localhost:8000", 0)

	license := handlers.IssueLicense(context.Background(), []int{})
	log.Println("got license, attempts: ", license.DigAllowed)

	treasures := handlers.Dig(context.Background(), server.DigParams{
		LicenseID: license.ID,
		PosX:      0,
		PosY:      0,
		Depth:     1,
	})
	log.Println("got treasures: ", treasures)

	for _, t := range treasures {
		coins := handlers.Cash(context.Background(), t)
		log.Println("got coins: ", coins)
	}
}
