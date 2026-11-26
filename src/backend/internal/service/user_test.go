package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/TwiLightDM/lab-track/internal/entity"
	errs "github.com/TwiLightDM/lab-track/internal/errors"
	"github.com/TwiLightDM/lab-track/internal/mocks"
	"github.com/TwiLightDM/lab-track/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func newUserServiceMocks(t *testing.T) (
	*mocks.MockUserRepository,
	*mocks.MockJWTGenerator,
	*service.UserService,
) {
	t.Helper()

	repo := mocks.NewMockUserRepository(t)
	jwt := mocks.NewMockJWTGenerator(t)

	svc := service.NewUserService(repo, jwt)

	return repo, jwt, svc
}

func requireServiceError(
	t *testing.T,
	err error,
	wantErr error,
	wantErrMsg string,
) {
	t.Helper()

	switch {
	case wantErr != nil:
		require.ErrorIs(t, err, wantErr)
	case wantErrMsg != "":
		require.EqualError(t, err, wantErrMsg)
	default:
		require.NoError(t, err)
	}
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name       string
		repoErr    error
		accessErr  error
		refreshErr error
		wantErr    error
		wantErrMsg string
	}{
		{
			name: "success",
		},
		{
			name:       "repository error",
			repoErr:    errors.New("database error"),
			wantErrMsg: "database error",
		},
		{
			name:    "duplicated email",
			repoErr: errs.ErrDuplicatedKey,
			wantErr: errs.ErrDuplicatedKey,
		},
		{
			name:       "access token error",
			accessErr:  errors.New("access token error"),
			wantErrMsg: "access token error",
		},
		{
			name:       "refresh token error",
			refreshErr: errors.New("refresh token error"),
			wantErrMsg: "refresh token error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, jwt, svc := newUserServiceMocks(t)

			repo.EXPECT().
				CreateUser(mock.Anything, mock.Anything).
				Run(func(_ context.Context, user *entity.User) {
					user.Id = 1

					require.Equal(t, "Ivan", user.Name)
					require.Equal(t, "Ivanov", user.Surname)
					require.Equal(t, "ivan@example.com", user.Email)
					require.Equal(t, "user", user.Role)
					require.True(t, user.IsActive)

					require.NoError(t, bcrypt.CompareHashAndPassword(
						[]byte(user.Password),
						[]byte("password"),
					))
				}).
				Return(tt.repoErr).
				Once()

			if tt.repoErr == nil {
				jwt.EXPECT().
					GenerateAccessToken(mock.Anything).
					Return("access-token", tt.accessErr).
					Once()

				if tt.accessErr == nil {
					jwt.EXPECT().
						GenerateRefreshToken(mock.Anything).
						Return("refresh-token", tt.refreshErr).
						Once()
				}
			}

			user, access, refresh, err := svc.CreateUser(
				context.Background(),
				"Ivan",
				"Ivanov",
				"ivan@example.com",
				"password",
				"user",
			)

			requireServiceError(t, err, tt.wantErr, tt.wantErrMsg)

			if tt.wantErr != nil || tt.wantErrMsg != "" {
				require.Nil(t, user)
				require.Empty(t, access)
				require.Empty(t, refresh)
				return
			}

			require.NotNil(t, user)
			require.Equal(t, 1, user.Id)
			require.Equal(t, "access-token", access)
			require.Equal(t, "refresh-token", refresh)
		})
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	expectedUsers := []entity.User{
		{Id: 1, Name: "Ivan"},
		{Id: 2, Name: "Petr"},
	}

	tests := []struct {
		name       string
		users      []entity.User
		repoErr    error
		wantErr    error
		wantErrMsg string
	}{
		{
			name:  "success",
			users: expectedUsers,
		},
		{
			name:       "repository error",
			repoErr:    errors.New("database error"),
			wantErrMsg: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, _, svc := newUserServiceMocks(t)

			repo.EXPECT().
				GetAllUsers(mock.Anything).
				Return(tt.users, tt.repoErr).
				Once()

			users, err := svc.GetAllUsers(context.Background())

			requireServiceError(t, err, tt.wantErr, tt.wantErrMsg)

			if tt.wantErr != nil || tt.wantErrMsg != "" {
				return
			}

			require.Equal(t, tt.users, users)
		})
	}
}

func TestUserService_GetUserById(t *testing.T) {
	expectedUser := &entity.User{
		Id:   1,
		Name: "Ivan",
	}

	tests := []struct {
		name    string
		user    *entity.User
		repoErr error
		wantErr error
	}{
		{
			name: "success",
			user: expectedUser,
		},
		{
			name:    "user not found",
			repoErr: errs.ErrEntityNotFound,
			wantErr: errs.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, _, svc := newUserServiceMocks(t)

			repo.EXPECT().
				GetUserById(mock.Anything, 1).
				Return(tt.user, tt.repoErr).
				Once()

			user, err := svc.GetUserById(context.Background(), 1)

			requireServiceError(t, err, tt.wantErr, "")

			if tt.wantErr != nil {
				require.Nil(t, user)
				return
			}

			require.Equal(t, expectedUser, user)
		})
	}
}

func TestUserService_PatchUser(t *testing.T) {
	name := "Petr"
	tests := []struct {
		name       string
		userName   *string
		surname    *string
		email      *string
		password   *string
		role       *string
		repoErr    error
		wantErr    error
		wantErrMsg string
		wantFields int
	}{
		{
			name:    "no fields provided",
			wantErr: errs.ErrNoFieldsProvided,
		},
		{
			name:       "success",
			userName:   &name,
			email:      new("petr@example.com"),
			password:   new("new-password"),
			role:       new("admin"),
			wantFields: 4,
		},
		{
			name:       "repository error",
			userName:   &name,
			repoErr:    errors.New("database error"),
			wantErrMsg: "database error",
			wantFields: 1,
		},
		{
			name:       "user not found",
			userName:   &name,
			repoErr:    errs.ErrEntityNotFound,
			wantErr:    errs.ErrEntityNotFound,
			wantFields: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, _, svc := newUserServiceMocks(t)

			if tt.wantFields > 0 {
				repo.EXPECT().
					PatchUser(
						mock.Anything,
						1,
						mock.MatchedBy(func(updates map[string]any) bool {
							if len(updates) != tt.wantFields {
								return false
							}

							if tt.userName != nil && updates["name"] != tt.userName {
								return false
							}

							if tt.surname != nil && updates["surname"] != tt.surname {
								return false
							}

							if tt.email != nil && updates["email"] != tt.email {
								return false
							}

							if tt.role != nil && updates["role"] != tt.role {
								return false
							}

							if tt.password != nil {
								hashedPassword, ok := updates["password"].(string)
								if !ok {
									return false
								}

								if bcrypt.CompareHashAndPassword(
									[]byte(hashedPassword),
									[]byte(*tt.password),
								) != nil {
									return false
								}
							}

							return true
						}),
					).
					Return(&entity.User{Id: 1, Name: "Petr"}, tt.repoErr).
					Once()
			}

			user, err := svc.PatchUser(
				context.Background(),
				1,
				tt.userName,
				tt.surname,
				tt.email,
				tt.password,
				tt.role,
			)

			requireServiceError(t, err, tt.wantErr, tt.wantErrMsg)

			if tt.wantErr != nil || tt.wantErrMsg != "" {
				require.Nil(t, user)
				return
			}

			require.NotNil(t, user)
			require.Equal(t, 1, user.Id)
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		repoErr    error
		wantErr    error
		wantErrMsg string
	}{
		{
			name: "success",
			id:   1,
		},
		{
			name:    "zero id",
			id:      0,
			wantErr: errs.ErrInvalidValue,
		},
		{
			name:       "repository error",
			id:         1,
			repoErr:    errors.New("database error"),
			wantErrMsg: "database error",
		},
		{
			name: "negative id is passed to repository",
			id:   -1,
		},
		{
			name:    "user not found",
			id:      1,
			repoErr: errs.ErrEntityNotFound,
			wantErr: errs.ErrEntityNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, _, svc := newUserServiceMocks(t)

			if tt.id != 0 {
				repo.EXPECT().
					DeleteUser(mock.Anything, tt.id).
					Return(tt.repoErr).
					Once()
			}

			err := svc.DeleteUser(context.Background(), tt.id)

			requireServiceError(t, err, tt.wantErr, tt.wantErrMsg)
		})
	}
}

func TestUserService_RefreshTokens(t *testing.T) {
	user := &entity.User{
		Id:   1,
		Role: "user",
	}

	tests := []struct {
		name        string
		repoErr     error
		accessErr   error
		refreshErr  error
		wantUser    *entity.User
		wantAccess  string
		wantRefresh string
		wantErr     error
	}{
		{
			name:        "success",
			wantUser:    user,
			wantAccess:  "access-token",
			wantRefresh: "refresh-token",
		},
		{
			name:    "user not found",
			repoErr: errs.ErrEntityNotFound,
			wantErr: errs.ErrEntityNotFound,
		},
		{
			name:      "access token error",
			accessErr: errors.New("access token error"),
			wantErr:   errors.New("access token error"),
		},
		{
			name:       "refresh token error",
			refreshErr: errors.New("refresh token error"),
			wantErr:    errors.New("refresh token error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, jwt, svc := newUserServiceMocks(t)

			repo.EXPECT().
				GetUserById(mock.Anything, 1).
				Return(user, tt.repoErr).
				Once()

			if tt.repoErr == nil {
				jwt.EXPECT().
					GenerateAccessToken(mock.Anything).
					Return("access-token", tt.accessErr).
					Once()

				if tt.accessErr == nil {
					jwt.EXPECT().
						GenerateRefreshToken(mock.Anything).
						Return("refresh-token", tt.refreshErr).
						Once()
				}
			}

			gotUser, gotAccess, gotRefresh, err := svc.RefreshTokens(
				context.Background(),
				1,
			)

			if tt.wantErr == nil {
				require.NoError(t, err)
			} else if errors.Is(tt.wantErr, errs.ErrEntityNotFound) {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.EqualError(t, err, tt.wantErr.Error())
			}

			require.Equal(t, tt.wantUser, gotUser)
			require.Equal(t, tt.wantAccess, gotAccess)
			require.Equal(t, tt.wantRefresh, gotRefresh)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	password := "password123"

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	require.NoError(t, err)

	user := &entity.User{
		Id:       1,
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     "user",
	}

	tests := []struct {
		name        string
		password    string
		repoErr     error
		accessErr   error
		refreshErr  error
		wantUser    *entity.User
		wantAccess  string
		wantRefresh string
		wantErr     error
	}{
		{
			name:        "success",
			password:    password,
			wantUser:    user,
			wantAccess:  "access-token",
			wantRefresh: "refresh-token",
		},
		{
			name:     "user not found",
			password: password,
			repoErr:  errs.ErrEntityNotFound,
			wantErr:  errs.ErrEntityNotFound,
		},
		{
			name:     "invalid password",
			password: "wrong-password",
			wantErr:  bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:      "access token error",
			password:  password,
			accessErr: errors.New("access token error"),
			wantErr:   errors.New("access token error"),
		},
		{
			name:       "refresh token error",
			password:   password,
			refreshErr: errors.New("refresh token error"),
			wantErr:    errors.New("refresh token error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, jwt, svc := newUserServiceMocks(t)

			repo.EXPECT().
				GetUserByEmail(mock.Anything, user.Email).
				Return(user, tt.repoErr).
				Once()

			credentialsInvalid := tt.repoErr != nil ||
				bcrypt.CompareHashAndPassword(
					[]byte(user.Password),
					[]byte(tt.password),
				) != nil

			if !credentialsInvalid {
				jwt.EXPECT().
					GenerateAccessToken(mock.Anything).
					Return("access-token", tt.accessErr).
					Once()

				if tt.accessErr == nil {
					jwt.EXPECT().
						GenerateRefreshToken(mock.Anything).
						Return("refresh-token", tt.refreshErr).
						Once()
				}
			}

			gotUser, gotAccess, gotRefresh, err := svc.Login(
				context.Background(),
				user.Email,
				tt.password,
			)

			if tt.wantErr == nil {
				require.NoError(t, err)
			} else if errors.Is(tt.wantErr, errs.ErrEntityNotFound) ||
				errors.Is(tt.wantErr, bcrypt.ErrMismatchedHashAndPassword) {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.EqualError(t, err, tt.wantErr.Error())
			}

			require.Equal(t, tt.wantUser, gotUser)
			require.Equal(t, tt.wantAccess, gotAccess)
			require.Equal(t, tt.wantRefresh, gotRefresh)
		})
	}
}
