package main

import (
	"context"
	"log"
	"time"

	"gold-rush/internal/server"
	"gold-rush/internal/service"
)

func main() {
	handlers := server.Init("http://localhost:8000", time.Millisecond*200)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()

		attempts := 0

		for attempts < 5 {
			if err := handlers.Healthcheck(); err != nil {
				log.Println("healthcheck: ", err)
				attempts++
			}

			<-time.After(time.Second)
		}
	}()

	service.Run(ctx, handlers)
}
