package service

import (
	"context"
	"fmt"

	"github.com/TwiLightDM/lab-track/internal/entity"
	errs "github.com/TwiLightDM/lab-track/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserById(ctx context.Context, id int) (*entity.User, error)
	PatchUser(ctx context.Context, id int, updates map[string]any) (*entity.User, error)
	DeleteUser(ctx context.Context, id int) error

	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type JWTGenerator interface {
	GenerateAccessToken(data map[string]any) (string, error)
	GenerateRefreshToken(data map[string]any) (string, error)
}

type UserService struct {
	repository   UserRepository
	jwtGenerator JWTGenerator
}

func NewUserService(repository UserRepository, jwtGenerator JWTGenerator) *UserService {
	return &UserService{repository: repository, jwtGenerator: jwtGenerator}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	name,
	surname,
	email,
	password,
	role string,
) (*entity.User, string, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	user := entity.User{
		Name:     name,
		Surname:  surname,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		IsActive: true,
	}

	if err = s.repository.CreateUser(ctx, &user); err != nil {
		return nil, "", "", err
	}

	data := make(map[string]any)
	data["id"] = user.Id
	data["role"] = user.Role

	accessToken, refreshToken, err := s.generateTokens(data)
	if err != nil {
		return nil, "", "", err
	}

	return &user, accessToken, refreshToken, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return s.repository.GetAllUsers(ctx)
}

func (s *UserService) GetUserById(ctx context.Context, id int) (*entity.User, error) {
	return s.repository.GetUserById(ctx, id)
}

func (s *UserService) PatchUser(
	ctx context.Context,
	id int,
	name,
	surname,
	email,
	password,
	role *string,
) (*entity.User, error) {
	if name == nil && surname == nil && email == nil && password == nil && role == nil {
		return nil, fmt.Errorf("нечего изменять: %w", errs.ErrNoFieldsProvided)
	}

	updates := make(map[string]any)

	if name != nil {
		updates["name"] = name
	}

	if surname != nil {
		updates["surname"] = surname
	}

	if email != nil {
		updates["email"] = email
	}

	if password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		updates["password"] = string(hashedPassword)
	}

	if role != nil {
		updates["role"] = role
	}

	user, err := s.repository.PatchUser(ctx, id, updates)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("id не может быть 0: %w", errs.ErrInvalidValue)
	}

	if err := s.repository.DeleteUser(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *UserService) RefreshTokens(ctx context.Context, id int) (*entity.User, string, string, error) {
	user, err := s.repository.GetUserById(ctx, id)
	if err != nil {
		return nil, "", "", err
	}

	data := make(map[string]any)
	data["id"] = user.Id
	data["role"] = user.Role

	accessToken, refreshToken, err := s.generateTokens(data)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*entity.User, string, string, error) {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", "", err
	}

	data := make(map[string]any)
	data["id"] = user.Id
	data["role"] = user.Role

	accessToken, refreshToken, err := s.generateTokens(data)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) generateTokens(data map[string]any) (string, string, error) {
	accessToken, err := s.jwtGenerator.GenerateAccessToken(data)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.jwtGenerator.GenerateRefreshToken(data)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
