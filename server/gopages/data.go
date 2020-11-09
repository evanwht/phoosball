package gopages

// Player from the db
type Player struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	NickName string `json:"display_name"`
}

type gameDBData struct {
	ID      int    `json:"id"`
	Date    string `json:"game_date"`
	T1pd    string `json:"team_1_d"`
	T1pdID  string `json:"team_1_d_id"`
	T1po    string `json:"team_1_o"`
	T1poID  string `json:"team_1_o_id"`
	T2pd    string `json:"team_2_d"`
	T2pdID  string `json:"team_2_d_id"`
	T2po    string `json:"team_2_o"`
	T2poID  string `json:"team_2_o_id"`
	T1half  int    `json:"team_1_half"`
	T2half  int    `json:"team_2_half"`
	T1final int    `json:"team_1_final"`
	T2final int    `json:"team_2_final"`
}

type gameData struct {
	gameDBData
	Team1Class string
	Team2Class string
}
