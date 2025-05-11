package types

import "github.com/rotisserie/eris"

var (
	ErrDatabaseQuery = eris.New("database query error")
	ErrNoRows        = eris.New("no rows found")
)
