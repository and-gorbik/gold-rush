package app

import "gold-rush/models"

// очередь с приоритетом
type AreaQueue struct {
	top    *models.ExploredArea
	length int
}

type Point struct {
	X int
	Y int
}

func (a *AreaQueue) Push(elem models.ExploredArea) {

}

func (a *AreaQueue) Pop() models.ExploredArea {
	if a.top == nil {
		return models.ExploredArea{}
	}

	return *a.top
}

func (a *AreaQueue) Len() int {
	return a.length
}
