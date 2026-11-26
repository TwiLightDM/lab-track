package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/TwiLightDM/lab-track/internal/dto"
	"github.com/TwiLightDM/lab-track/internal/entity"
	errs "github.com/TwiLightDM/lab-track/internal/errors"
	"github.com/labstack/echo/v4"
)

type userService interface {
	CreateUser(ctx context.Context, name, surname, email, password, role string) (*entity.User, string, string, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetUserById(ctx context.Context, id int) (*entity.User, error)
	PatchUser(ctx context.Context, id int, name, surname, email, password, role *string) (*entity.User, error)
	DeleteUser(ctx context.Context, id int) error

	Login(ctx context.Context, email, password string) (*entity.User, string, string, error)
	RefreshTokens(ctx context.Context, id int) (*entity.User, string, string, error)
}

type validatorService interface {
	ValidateError(c echo.Context, err error) error
	ValidateRequest(s any) error
}

type UserHandler struct {
	service   userService
	validator validatorService
}

func NewUserHandler(service userService, tools validatorService) *UserHandler {
	return &UserHandler{service: service, validator: tools}
}

// CreateUser
// @Summary      Создать пользователя
// @Description  Создаёт пользователя и возвращает пару токенов.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UserRequest  true  "Данные нового пользователя"
// @Success      201      {object}  dto.UserResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      401      {object}  dto.ErrorResponse
// @Failure      403      {object}  dto.ErrorResponse
// @Failure      409      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var request dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return h.validator.ValidateError(c, errs.ErrInvalidJSON)
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Surname = strings.TrimSpace(request.Surname)
	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)
	request.Role = strings.TrimSpace(request.Role)

	if err := h.validator.ValidateRequest(request); err != nil {
		return h.validator.ValidateError(c, err)
	}

	user, accessToken, refreshToken, err := h.service.CreateUser(
		c.Request().Context(),
		request.Name,
		request.Surname,
		request.Email,
		request.Password,
		request.Role,
	)
	if err != nil {
		return h.validator.ValidateError(c, err)
	}

	return c.JSON(http.StatusCreated, dto.UserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
		Tokens: &dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// GetAllUsers
// @Summary      Список пользователей
// @Tags         users
// @Produce      json
// @Success      200  {array}   dto.UserResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := h.service.GetAllUsers(c.Request().Context())
	if err != nil {
		return h.validator.ValidateError(c, err)
	}

	response := make([]dto.UserResponse, 0, len(users))

	for _, user := range users {
		response = append(response, dto.UserResponse{
			Id:       user.Id,
			Name:     user.Name,
			Surname:  user.Surname,
			Email:    user.Email,
			Role:     user.Role,
			IsActive: user.IsActive,
		})
	}

	return c.JSON(http.StatusOK, response)
}

// GetUserById
// @Summary      Получить пользователя
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.validator.ValidateError(c, fmt.Errorf("невалидный id: %w", errs.ErrInvalidValue))
	}

	user, err := h.service.GetUserById(c.Request().Context(), id)
	if err != nil {
		return h.validator.ValidateError(c, err)
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	})
}

// PatchUser
// @Summary      Обновить пользователя
// @Description  Обновляет только переданные поля.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "ID пользователя"
// @Param        request  body      dto.PatchUserRequest   true  "Поля для обновления"
// @Success      200      {object}  dto.UserResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [patch]
func (h *UserHandler) PatchUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.validator.ValidateError(c, fmt.Errorf("невалидный id: %w", errs.ErrInvalidValue))
	}

	var request dto.PatchUserRequest
	if err = c.Bind(&request); err != nil {
		return h.validator.ValidateError(c, errs.ErrInvalidJSON)
	}

	if request.Name != nil {
		request.Name = new(strings.TrimSpace(*request.Name))
	}
	if request.Surname != nil {
		request.Surname = new(strings.TrimSpace(*request.Surname))
	}
	if request.Email != nil {
		request.Email = new(strings.TrimSpace(*request.Email))
	}
	if request.Password != nil {
		request.Password = new(strings.TrimSpace(*request.Password))
	}
	if request.Role != nil {
		request.Role = new(strings.TrimSpace(*request.Role))
	}

	if err = h.validator.ValidateRequest(request); err != nil {
		return h.validator.ValidateError(c, err)
	}

	user, err := h.service.PatchUser(
		c.Request().Context(),
		id,
		request.Name,
		request.Surname,
		request.Email,
		request.Password,
		request.Role,
	)
	if err != nil {
		return h.validator.ValidateError(c, err)
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	})
}

// DeleteUser
// @Summary      Удалить пользователя
// @Tags         users
// @Param        id   path  int  true  "ID пользователя"
// @Success      204  "Пользователь удалён"
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return h.validator.ValidateError(c, fmt.Errorf("невалидный id: %w", errs.ErrInvalidValue))
	}

	err = h.service.DeleteUser(c.Request().Context(), id)
	if err != nil {
		return h.validator.ValidateError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetMe
// @Summary      Получить профиль текущего пользователя
// @Description  Возвращает данные авторизованного пользователя на основе userId из JWT-токена
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.UserResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /users/me [get]
func (h *UserHandler) GetMe(c echo.Context) error {
	userId := c.Get("userId")
	id := int(userId.(float64))

	user, err := h.service.GetUserById(c.Request().Context(), id)
	if err != nil {
		return h.validator.ValidateError(c, err)
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	})
}

// RefreshTokens
// @Summary      Обновить токены
// @Description  Требует действующий refresh-токен в заголовке Authorization.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dto.UserResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /public/refresh [post]
func (h *UserHandler) RefreshTokens(c echo.Context) error {
	userId := c.Get("userId")
	id := int(userId.(float64))

	user, accessToken, refreshToken, err := h.service.RefreshTokens(c.Request().Context(), id)
	if err != nil {
		return h.validator.ValidateError(c, err)
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
		Tokens: &dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// Login
// @Summary      Вход
// @Description  Возвращает access/refresh токены.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "Email и пароль"
// @Success      200      {object}  dto.UserResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      401      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /public/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var request dto.LoginRequest
	if err := c.Bind(&request); err != nil {
		return h.validator.ValidateError(c, errs.ErrInvalidJSON)
	}

	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)

	if err := h.validator.ValidateRequest(request); err != nil {
		return h.validator.ValidateError(c, err)
	}

	user, accessToken, refreshToken, err := h.service.Login(c.Request().Context(), request.Email, request.Password)
	if err != nil {
		return h.validator.ValidateError(c, err)
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
		Tokens: &dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}
