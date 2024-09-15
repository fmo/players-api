package domain

type Player struct {
	Age                 *int    `json:"age,omitempty"`
	Id                  *string `json:"id,omitempty"`
	Name                *string `json:"name,omitempty"`
	Position            *string `json:"position,omitempty"`
	ShirtNumber         *string `json:"shirtNumber,omitempty"`
	TeamId              *int    `json:"teamId,omitempty"`
	MarketValue         int     `json:"marketValue"`
	MarketValueCurrency string  `json:"marketValueCurrency"`
}
