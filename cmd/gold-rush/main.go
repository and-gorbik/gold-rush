package main

import (
	"gold-rush/internal/server"
	"gold-rush/internal/service"
)

func main() {
	handlers := server.Init("http://localhost:8000")
	service.Run(handlers)
}
