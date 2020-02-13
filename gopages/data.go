package gopages

// Player from the db
type Player struct {
	ID       int
	Name     string
	NickName string
}

type gameDBData struct {
	ID         int    `json:"id"`
	Date       string `json:"game_date"`
	T1pd       string `json:"team_1_Defense"`
	T1po       string `json:"team_1_Offense"`
	T2pd       string `json:"team_2_Defense"`
	T2po       string `json:"team_2_Offense"`
	T1half     int    `json:"team_1_half"`
	T2half     int    `json:"team_2_half"`
	T1final    int    `json:"team_1_final"`
	T2final    int    `json:"team_2_final"`
}

type gameData struct {
	gameDBData
	Team1Class string
	Team2Class string
}