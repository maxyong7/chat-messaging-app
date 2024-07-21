package entity

import "database/sql"

type NullString struct {
	sql.NullString
}

func (s NullString) String() string {
	return s.NullString.String
}
