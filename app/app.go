package app

import (
	"log"
	"net/http"
	"time"

	"gold-rush/models"
	"gold-rush/server"
)

const (
	maxArea             = 3500
	maxDepth            = 10
	totalTreasuresCount = 490000
	totalCoinsCount     = 23030000
)

type App struct {
	provider goldRushServer
}

func New() *App {
	return &App{
		provider: server.GoldRushServer{
			client: &http.Client{
				Timeout: 5 * time.Second,
			}
		},
	}
}

func (a *App) Run() {
	area := models.Area{
		PosX:  0,
		PosY:  0,
		SizeX: 10,
		SizeY: 10,
	}
	exploredArea, err := a.provider.Explore(area)
	if err != nil {
		message, ok := readBusinessError(err)
		if !ok {
			log.Printf("Server error: %s\n", message)
		} else {
			log.Printf("Client error: %s\n", message)
		}
	}
}
