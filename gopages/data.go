package gopages

// Player from the db
type Player struct {
	ID       int
	Name     string
	NickName string
}

type gameData struct {
	ID         int
	Date       string
	T1pd       string
	T1po       string
	T2pd       string
	T2po       string
	T1half     int
	T2half     int
	T1final    int
	T2final    int
	Team1Class string
	Team2Class string
}