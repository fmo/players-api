package models

type Player struct {
	Team        string `json:"team"`
	TeamId      int    `json:"teamId"`
	Name        string `json:"name"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
	Photo       string `json:"photo"`
	RapidApiID  string `json:"id"`
	Appearances int    `json:"appearances"`
	Position    string `json:"position"`
}
