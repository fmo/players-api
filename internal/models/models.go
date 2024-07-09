package models

type Player struct {
	Team        string `json:"team"`
	TeamId      int32  `json:"teamId"`
	Name        string `json:"name"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Age         int32  `json:"age"`
	Nationality string `json:"nationality"`
	Photo       string `json:"photo"`
	RapidApiID  string `json:"id"`
	Appearances int32  `json:"appearances"`
	Position    string `json:"position"`
}
