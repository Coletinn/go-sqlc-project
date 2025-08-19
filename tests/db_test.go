package tests

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"  // Adicione este import
	"github.com/stretchr/testify/require"
)

// SetupTestDB é uma função helper, não um teste
func SetupTestDB(t *testing.T) *sql.DB {
	// Use um banco de dados de teste diferente
	testDBURL := "postgres://gustavo:1910@localhost:5432/postgres?sslmode=disable"
	
	conn, err := sql.Open("postgres", testDBURL)
	require.NoError(t, err, "Failed to connect to test database")
	
	// Verificar conexão
	err = conn.Ping()
	require.NoError(t, err, "Failed to ping test database")
	
	// Limpar dados existentes (opcional, dependendo da estratégia)
	_, err = conn.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Logf("Warning: Failed to truncate users table: %v", err)
	}
	
	return conn
}