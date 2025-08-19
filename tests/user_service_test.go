package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sqlc-testing/db"
	"sqlc-testing/services"
)

func TestUserService_CreateUser(t *testing.T) {
	conn := SetupTestDB(t)
	defer conn.Close()
	
	userService := service.NewUserService(conn)
	ctx := context.Background()

	tests := []struct {
		name        string
		userName    string
		email       string
		phone       string
		expectError bool
	}{
		{
			name:        "valid user with phone",
			userName:    "João Silva",
			email:       "joao@email.com",
			phone:       "11999999999",
			expectError: false,
		},
		{
			name:        "valid user without phone",
			userName:    "Maria Santos",
			email:       "maria@email.com",
			phone:       "",
			expectError: false,
		},
		{
			name:        "duplicate email",
			userName:    "Pedro Costa",
			email:       "joao@email.com", // email já usado no primeiro teste
			phone:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userService.CreateUser(ctx, tt.userName, tt.email, tt.phone)
			
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Greater(t, user.ID, int32(0))
			assert.Equal(t, tt.userName, user.Name)
			assert.Equal(t, tt.email, user.Email)
			
			if tt.phone != "" {
				assert.True(t, user.Phone.Valid)
				assert.Equal(t, tt.phone, user.Phone.String)
			} else {
				assert.False(t, user.Phone.Valid)
			}
			
			assert.WithinDuration(t, time.Now(), user.CreatedAt.Time, 5*time.Second)
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	conn := SetupTestDB(t)
	defer conn.Close()
	
	userService := service.NewUserService(conn)
	ctx := context.Background()

	// Criar um usuário para teste
	createdUser, err := userService.CreateUser(ctx, "Test User", "test@email.com", "11888888888")
	require.NoError(t, err)

	t.Run("existing user", func(t *testing.T) {
		user, err := userService.GetUserByID(ctx, createdUser.ID)
		require.NoError(t, err)
		
		assert.Equal(t, createdUser.ID, user.ID)
		assert.Equal(t, createdUser.Name, user.Name)
		assert.Equal(t, createdUser.Email, user.Email)
		assert.Equal(t, createdUser.Phone.String, user.Phone.String)
	})

	t.Run("non-existing user", func(t *testing.T) {
		_, err := userService.GetUserByID(ctx, 99999)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}

func TestUserService_ListUsers(t *testing.T) {
	conn := SetupTestDB(t)
	defer conn.Close()
	
	userService := service.NewUserService(conn)
	ctx := context.Background()

	// Criar alguns usuários
	users := []struct {
		name  string
		email string
		phone string
	}{
		{"User 1", "user1@email.com", "11111111111"},
		{"User 2", "user2@email.com", ""},
		{"User 3", "user3@email.com", "33333333333"},
	}

	var createdUsers []db.User
	for _, u := range users {
		user, err := userService.CreateUser(ctx, u.name, u.email, u.phone)
		require.NoError(t, err)
		createdUsers = append(createdUsers, user)
	}

	// Testar listagem
	listedUsers, err := userService.ListUsers(ctx)
	require.NoError(t, err)
	
	assert.Len(t, listedUsers, len(users))
	
	// Verificar se todos os usuários criados estão na lista
	for _, created := range createdUsers {
		found := false
		for _, listed := range listedUsers {
			if listed.ID == created.ID {
				assert.Equal(t, created.Name, listed.Name)
				assert.Equal(t, created.Email, listed.Email)
				found = true
				break
			}
		}
		assert.True(t, found, "Created user not found in list: %v", created)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	conn := SetupTestDB(t)
	defer conn.Close()
	
	userService := service.NewUserService(conn)
	ctx := context.Background()

	// Criar usuário para atualizar
	originalUser, err := userService.CreateUser(ctx, "Original Name", "original@email.com", "11999999999")
	require.NoError(t, err)

	t.Run("update all fields", func(t *testing.T) {
		updatedUser, err := userService.UpdateUser(ctx, originalUser.ID, "Updated Name", "updated@email.com", "11888888888")
		require.NoError(t, err)
		
		assert.Equal(t, originalUser.ID, updatedUser.ID)
		assert.Equal(t, "Updated Name", updatedUser.Name)
		assert.Equal(t, "updated@email.com", updatedUser.Email)
		assert.True(t, updatedUser.Phone.Valid)
		assert.Equal(t, "11888888888", updatedUser.Phone.String)
	})

	t.Run("update without phone", func(t *testing.T) {
		updatedUser, err := userService.UpdateUser(ctx, originalUser.ID, "No Phone User", "nophone@email.com", "")
		require.NoError(t, err)
		
		assert.Equal(t, "No Phone User", updatedUser.Name)
		assert.Equal(t, "nophone@email.com", updatedUser.Email)
		assert.False(t, updatedUser.Phone.Valid)
	})

	t.Run("update non-existing user", func(t *testing.T) {
		_, err := userService.UpdateUser(ctx, 99999, "Fake User", "fake@email.com", "")
		assert.Error(t, err)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	conn := SetupTestDB(t)
	defer conn.Close()
	
	userService := service.NewUserService(conn)
	ctx := context.Background()

	// Criar usuário para deletar
	user, err := userService.CreateUser(ctx, "To Delete", "delete@email.com", "")
	require.NoError(t, err)

	t.Run("delete existing user", func(t *testing.T) {
		err := userService.DeleteUser(ctx, user.ID)
		require.NoError(t, err)
		
		// Verificar se foi deletado
		_, err = userService.GetUserByID(ctx, user.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("delete non-existing user", func(t *testing.T) {
		err := userService.DeleteUser(ctx, 99999)
		// DeleteUser pode não retornar erro mesmo se o usuário não existir
		// dependendo da implementação do SQLC. Teste conforme sua implementação.
		assert.NoError(t, err) // ou assert.Error(t, err) se sua implementação retorna erro
	})
}

// Benchmark para testar performance
func BenchmarkUserService_CreateUser(b *testing.B) {
	conn := SetupTestDB(&testing.T{}) // Hack para usar SetupTestDB
	defer conn.Close()
	
	userService := service.NewUserService(conn)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := userService.CreateUser(ctx, "Benchmark User", "bench@email.com", "11999999999")
		if err != nil {
			b.Fatal(err)
		}
		
		// Limpar para próxima iteração
		conn.Exec("TRUNCATE TABLE users RESTART IDENTITY")
	}
}

// Helper para testes de integração
func TestUserService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	conn := SetupTestDB(t)
	defer conn.Close()
	
	userService := service.NewUserService(conn)
	ctx := context.Background()

	// Teste completo: Create -> Read -> Update -> Delete
	t.Run("full CRUD cycle", func(t *testing.T) {
		// Create
		user, err := userService.CreateUser(ctx, "Integration User", "integration@email.com", "11777777777")
		require.NoError(t, err)
		userID := user.ID

		// Read
		fetchedUser, err := userService.GetUserByID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, user.Name, fetchedUser.Name)

		// Update
		updatedUser, err := userService.UpdateUser(ctx, userID, "Updated Integration", "updated-integration@email.com", "11666666666")
		require.NoError(t, err)
		assert.Equal(t, "Updated Integration", updatedUser.Name)

		// Delete
		err = userService.DeleteUser(ctx, userID)
		require.NoError(t, err)

		// Verify deletion
		_, err = userService.GetUserByID(ctx, userID)
		assert.Error(t, err)
	})
}