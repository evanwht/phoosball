package players

// Player from the db
type Player struct {
	ID          int          `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	NickName    string       `json:"nickname,omitempty" db:"display_name"`
	Email       string       `json:"email,omitempty" db:"email"`
	Password    string       `json:"-" db:"password"`
	UserName    string       `json:"username" db:"username"`
	Permissions []Permission `json:"permissions" db:"permissions"`
}

type Permission struct {
	Type string `json:"type" db:"type"`
}
