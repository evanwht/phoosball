package util

import (
	"github.com/jmoiron/sqlx"
)

// Env : environment variables that should be shared between routes but created ony once
type Env struct {
	DB *sqlx.DB
}
