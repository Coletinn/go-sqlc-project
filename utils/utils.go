package utils

import "database/sql"

// Helper to convert string to sql.NullString
func NullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}