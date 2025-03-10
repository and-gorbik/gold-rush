package server

type License struct {
	ID         int `json:"id"`
	DigAllowed int `json:"digAllowed"`
	DigUsed    int `json:"digUsed"`
}

type Area struct {
	PosX  int `json:"posX"`
	PosY  int `json:"posY"`
	SizeX int `json:"sizeX"`
	SizeY int `json:"sizeY"`
}

type ExploredArea struct {
	Area   Area `json:"area"`
	Amount int  `json:"amount"`
}

type DigParams struct {
	LicenseID int `json:"licenseID"`
	PosX      int `json:"posX"`
	PosY      int `json:"posY"`
	Depth     int `json:"depth"`
}
