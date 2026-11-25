package tools

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/TwiLightDM/lab-track/internal/dto"
	errs "github.com/TwiLightDM/lab-track/internal/errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

const DefaultPage = 1
const DefaultLimit = 10

type ValidatorService struct {
	validator *validator.Validate
}

func NewValidatorService() *ValidatorService {
	return &ValidatorService{validator: validator.New()}
}

func (v *ValidatorService) ValidateError(c echo.Context, err error) error {
	var status int

	switch {
	case errors.Is(err, errs.ErrInvalidValue),
		errors.Is(err, errs.ErrNoFieldsProvided),
		errors.Is(err, errs.ErrInvalidJSON):
		status = http.StatusBadRequest
	case errors.Is(err, errs.ErrInvalidCredentials),
		errors.Is(err, errs.ErrMissingToken),
		errors.Is(err, errs.ErrInvalidToken),
		errors.Is(err, errs.ErrTokenExpired),
		errors.Is(err, errs.ErrUnknownRole):
		status = http.StatusUnauthorized
	case errors.Is(err, errs.ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(err, errs.ErrEntityNotFound):
		status = http.StatusNotFound
	case errors.Is(err, context.Canceled):
		status = http.StatusRequestTimeout
	case errors.Is(err, errs.ErrDuplicatedKey),
		errors.Is(err, errs.ErrForeignKeyViolated):
		status = http.StatusConflict
	case errors.Is(err, context.DeadlineExceeded):
		status = http.StatusGatewayTimeout
	default:
		status = http.StatusInternalServerError
	}

	return c.JSON(status, dto.ErrorResponse{Code: errorCode(err), Message: err.Error()})
}

func errorCode(err error) string {
	errorCodes := []struct {
		err  error
		code string
	}{
		{errs.ErrInvalidValue, "ErrInvalidValue"},
		{errs.ErrEntityNotFound, "ErrEntityNotFound"},
		{errs.ErrNoFieldsProvided, "ErrNoFieldsProvided"},
		{errs.ErrDuplicatedKey, "ErrDuplicatedKey"},
		{errs.ErrForeignKeyViolated, "ErrForeignKeyViolated"},
		{errs.ErrForbidden, "ErrForbidden"},
		{errs.ErrInvalidCredentials, "ErrInvalidCredentials"},
		{errs.ErrMissingToken, "ErrMissingToken"},
		{errs.ErrInvalidToken, "ErrInvalidToken"},
		{errs.ErrUnknownRole, "ErrUnknownRole"},
		{errs.ErrInvalidJSON, "ErrInvalidJSON"},
		{errs.ErrTokenExpired, "ErrTokenExpired"},
	}

	for _, item := range errorCodes {
		if errors.Is(err, item.err) {
			return item.code
		}
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return "ErrDeadlineExceeded"
	}

	if errors.Is(err, context.Canceled) {
		return "ErrContextCanceled"
	}

	return "ErrInternalServer"
}

func (v *ValidatorService) ValidateRequest(s any) error {
	if err := v.validator.Struct(s); err != nil {
		if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
			msgs := make([]string, 0, len(ve))
			for _, fe := range ve {
				msgs = append(msgs, fmt.Sprintf("%s: %s", fe.Field(), msgForTag(fe.Tag())))
			}
			return fmt.Errorf("%s: %w", strings.Join(msgs, "; "), errs.ErrInvalidValue)
		}
		return fmt.Errorf("ошибка валидации: %w", errs.ErrInvalidValue)
	}
	return nil
}

func (v *ValidatorService) GetPageAndLimit(c echo.Context) (int, int) {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil && page <= 0 {
		page = DefaultPage
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil && limit <= 0 {
		limit = DefaultLimit
	}

	return page, limit
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "обязательное поле"
	case "email":
		return "некорректный email"
	case "min":
		return "значение слишком короткое"
	case "oneof":
		return "недопустимое значение"
	default:
		return "невалидное значение"
	}
}
