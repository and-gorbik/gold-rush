package app

import (
	"gold-rush/models"
	"log"
)

const (
	maxArea             = 3500
	maxDepth            = 10
	totalTreasuresCount = 490000
	totalCoinsCount     = 23030000
)

type App struct {
	provider provider
}

func New(provider provider) *App {
	return &App{
		provider: provider,
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
