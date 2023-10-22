package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/infrastructure/auth"
)

func Guard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")
		if authorization == "" {
			return apperror.Unauthorized(apperror.Payload{
				Status: apperror.StatusErrorInvalidAuth,
			})
		}

		s := strings.Split(authorization, " ")
		if len(s) < 2 || s[0] != "Bearer" {
			return apperror.Unauthorized(apperror.Payload{
				Status: apperror.StatusErrorInvalidAuth,
			})
		}

		sessionClaim, err := auth.Get().VerifyToken(s[1])
		if err != nil {
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
