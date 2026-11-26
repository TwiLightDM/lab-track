package middlewares

import (
	"strings"

	errs "github.com/TwiLightDM/lab-track/internal/errors"
	"github.com/labstack/echo/v4"
)

type ParseJWTService interface {
	ParseJWT(token string) (map[string]any, error)
}

type ValidatorService interface {
	ValidateError(c echo.Context, err error) error
}

type MWManager struct {
	parseJWT  ParseJWTService
	validator ValidatorService
}

func NewMWManager(parseJWTService ParseJWTService, validator ValidatorService) *MWManager {
	return &MWManager{parseJWT: parseJWTService, validator: validator}
}

func (mw *MWManager) AdminCheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role")
			role, ok := userRole.(string)
			if !ok {
				return mw.validator.ValidateError(c, errs.ErrUnknownRole)
			}

			if role != "admin" {
				return mw.validator.ValidateError(c, errs.ErrForbidden)
			}

			return next(c)
		}
	}
}

func (mw *MWManager) BearerAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return mw.validator.ValidateError(c, errs.ErrMissingToken)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return mw.validator.ValidateError(c, errs.ErrInvalidToken)
			}

			claims, err := mw.parseJWT.ParseJWT(parts[1])
			if err != nil {
				return mw.validator.ValidateError(c, err)
			}

			c.Set("userId", claims["sub"])
			c.Set("role", claims["role"])

			return next(c)
		}
	}
}
