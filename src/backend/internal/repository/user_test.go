package repository_test

import (
	"context"
	"testing"

	"github.com/TwiLightDM/lab-track/internal/entity"
	errs "github.com/TwiLightDM/lab-track/internal/errors"
	"github.com/TwiLightDM/lab-track/internal/repository"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupUserRepository(t *testing.T) *repository.UserRepository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		TranslateError: true,
	})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)

	require.NoError(t, db.AutoMigrate(&entity.User{}))

	return repository.NewUserRepository(db)
}

func createTestUser(
	t *testing.T,
	repo *repository.UserRepository,
	email string,
) *entity.User {
	t.Helper()

	user := &entity.User{
		Name:     "Ivan",
		Surname:  "Ivanov",
		Email:    email,
		Password: "password",
		Role:     "user",
		IsActive: true,
	}

	require.NoError(t, repo.CreateUser(context.Background(), user))

	return user
}

func TestUserRepository_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*testing.T, *repository.UserRepository)
		user    *entity.User
		wantErr error
	}{
		{
			name: "success",
			user: &entity.User{
				Name:     "Ivan",
				Surname:  "Ivanov",
				Email:    "ivan@example.com",
				Password: "password",
				Role:     "user",
				IsActive: true,
			},
		},
		{
			name: "duplicated email",
			prepare: func(t *testing.T, repo *repository.UserRepository) {
				createTestUser(t, repo, "ivan@example.com")
			},
			user: &entity.User{
				Name:  "Petr",
				Email: "ivan@example.com",
			},
			wantErr: errs.ErrDuplicatedKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupUserRepository(t)

			if tt.prepare != nil {
				tt.prepare(t, repo)
			}

			err := repo.CreateUser(context.Background(), tt.user)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.NotZero(t, tt.user.Id)
		})
	}
}

func TestUserRepository_GetAllUsers(t *testing.T) {
	tests := []struct {
		name          string
		prepare       func(*testing.T, *repository.UserRepository)
		expectedCount int
	}{
		{
			name:          "empty list",
			expectedCount: 0,
		},
		{
			name: "users exist",
			prepare: func(t *testing.T, repo *repository.UserRepository) {
				createTestUser(t, repo, "ivan@example.com")
				createTestUser(t, repo, "petr@example.com")
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupUserRepository(t)

			if tt.prepare != nil {
				tt.prepare(t, repo)
			}

			users, err := repo.GetAllUsers(context.Background())

			require.NoError(t, err)
			require.Len(t, users, tt.expectedCount)
		})
	}
}

func TestUserRepository_GetUserById(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*testing.T, *repository.UserRepository) int
		wantErr error
	}{
		{
			name: "success",
			prepare: func(t *testing.T, repo *repository.UserRepository) int {
				return createTestUser(t, repo, "ivan@example.com").Id
			},
		},
		{
			name: "user not found",
			prepare: func(t *testing.T, repo *repository.UserRepository) int {
				return 999
			},
			wantErr: errs.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupUserRepository(t)
			id := tt.prepare(t, repo)

			user, err := repo.GetUserById(context.Background(), id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, id, user.Id)
		})
	}
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		prepare func(*testing.T, *repository.UserRepository)
		wantErr error
	}{
		{
			name:  "success",
			email: "ivan@example.com",
			prepare: func(t *testing.T, repo *repository.UserRepository) {
				createTestUser(t, repo, "ivan@example.com")
			},
		},
		{
			name:    "user not found",
			email:   "unknown@example.com",
			wantErr: errs.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupUserRepository(t)

			if tt.prepare != nil {
				tt.prepare(t, repo)
			}

			user, err := repo.GetUserByEmail(context.Background(), tt.email)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, tt.email, user.Email)
		})
	}
}

func TestUserRepository_PatchUser(t *testing.T) {
	tests := []struct {
		name     string
		prepare  func(*testing.T, *repository.UserRepository) int
		updates  map[string]any
		wantName string
		wantErr  error
	}{
		{
			name: "success",
			prepare: func(t *testing.T, repo *repository.UserRepository) int {
				return createTestUser(t, repo, "ivan@example.com").Id
			},
			updates:  map[string]any{"name": "Petr"},
			wantName: "Petr",
		},
		{
			name: "user not found",
			prepare: func(t *testing.T, repo *repository.UserRepository) int {
				return 999
			},
			updates: map[string]any{"name": "Petr"},
			wantErr: errs.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupUserRepository(t)
			id := tt.prepare(t, repo)

			user, err := repo.PatchUser(
				context.Background(),
				id,
				tt.updates,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, tt.wantName, user.Name)
		})
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*testing.T, *repository.UserRepository) int
		wantErr error
	}{
		{
			name: "success",
			prepare: func(t *testing.T, repo *repository.UserRepository) int {
				return createTestUser(t, repo, "ivan@example.com").Id
			},
		},
		{
			name: "user not found",
			prepare: func(t *testing.T, repo *repository.UserRepository) int {
				return 999
			},
			wantErr: errs.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupUserRepository(t)
			id := tt.prepare(t, repo)

			err := repo.DeleteUser(context.Background(), id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)

			_, err = repo.GetUserById(context.Background(), id)
			require.ErrorIs(t, err, errs.ErrEntityNotFound)
		})
	}
}
