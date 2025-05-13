package types

import "github.com/rotisserie/eris"

var (
	ErrDatabaseQuery = eris.New("database query error")
	ErrNoRows        = eris.New("no rows found")
)

type HTTPError struct {
	Code      int    `json:"code" example:"400"`
	Message   string `json:"message" example:"Invalid request payload"`
	Error     string `json:"error,omitempty" example:"Bad Request"`
	Timestamp string `json:"timestamp,omitempty" example:"2023-12-01T12:34:56Z"`
	Path      string `json:"path,omitempty" example:"/books"`
}
