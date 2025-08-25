package tests

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func SetupTestDB(t *testing.T) *sql.DB {
	testDBURL := os.Getenv("DATABASE_URL")
	
	conn, err := sql.Open("postgres", testDBURL)
	require.NoError(t, err, "Failed to connect to test database")
	
	err = conn.Ping()
	require.NoError(t, err, "Failed to ping test database")
	
	_, err = conn.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Logf("Warning: Failed to truncate users table: %v", err)
	}
	
	return conn
}
