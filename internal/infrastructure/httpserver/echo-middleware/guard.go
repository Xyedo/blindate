package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/labstack/echo/v4"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/infrastructure/auth"
)

func Guard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorizationHeader := c.Request().Header.Get("Authorization")
		if authorizationHeader == "" {
			return apperror.Unauthorized(apperror.Payload{
				Status: apperror.StatusErrorInvalidAuth,
			})
		}

		s := strings.Split(authorizationHeader, " ")
		if len(s) < 2 || s[0] != "Bearer" {
			return apperror.Unauthorized(apperror.Payload{
				Status: apperror.StatusErrorInvalidAuth,
			})
		}

		sessionClaim, err := auth.Get().VerifyToken(s[1])
		if err != nil {
			if errors.Is(err, jwt.ErrExpired) {
				return apperror.Unauthorized(apperror.Payload{
					Status: apperror.StatusErrorExpiredAuth,
					Error:  err,
				})
			}
			
			return apperror.Unauthorized(apperror.Payload{
				Status: apperror.StatusErrorInvalidAuth,
				Error:  err,
			})
		}
		ctx := context.WithValue(c.Request().Context(), auth.RequestId, sessionClaim.Claims.Subject)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
