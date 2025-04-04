package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/seemyown/backend-toolkit/btools/db"
	"github.com/seemyown/backend-toolkit/btools/exc"
	"strings"
)

var errLogger = log.NewSubLogger("error_middleware")

func ErrorMiddleware(locale string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()
		if err == nil {
			return nil
		}

		var appErr *exc.Error
		var fiberErr *fiber.Error
		var repositoryErr *db.RepositoryError
		if errors.As(err, &fiberErr) {
			appErr = exc.NewAppError(
				strings.TrimSpace(fiberErr.Message),
				toSnakeCase(fiberErr.Message),
				"",
				fiberErr.Code,
			)
		} else if errors.As(err, &repositoryErr) {
			appErr = exc.NewAppError(
				"InternalError",
				"repository_error",
				repositoryErr.MessageFor(locale),
				db.MapToHttpError[repositoryErr.Code],
			)
		} else if !errors.As(err, &appErr) {
			appErr = exc.InternalServerError("Unckown error")
		}

		errLogger.Error(err, "Request error %+v", appErr)
		return ctx.Status(appErr.StatusCode).JSON(appErr)
	}
}

func toSnakeCase(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "_"))
}
