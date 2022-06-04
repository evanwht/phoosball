package standings

// Standing : wins and losses of a player in a context
type Standing struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
}
