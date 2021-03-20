package models

type ServiceStatus map[string]interface{}

type Balance struct {
	Value  int   `json:"balance"`
	Wallet []int `json:"wallet"`
}

type License struct {
	ID         int `json:"id"`
	DigAllowed int `json:"digAllowed"`
}

type LicenseFull struct {
	License
	DigUsed int `json:"digUsed"`
}

type PaymentForLicense []int

type ExploredArea struct {
	Area   Area `json:"area"`
	Amount int  `json:"amount"`
}

type Area struct {
	PosX  int `json:"posX"`
	PosY  int `json:"posY"`
	SizeX int `json:"sizeX"`
	SizeY int `json:"sizeY"`
}

type DigParams struct {
	LicenseID int `json:"licenseID"`
	PosX      int `json:"posX"`
	PosY      int `json:"posY"`
	Depth     int `json:"depth"`
}

type TreasuresList []string

type Treasure string

type PaymentForTreasure []int
