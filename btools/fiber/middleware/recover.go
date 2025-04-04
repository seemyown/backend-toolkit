package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/seemyown/backend-toolkit/btools/exc"
)

var panicMiddlewareLogger = log.NewSubLogger("panic.handler")

func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				panicMiddlewareLogger.Trace("[PANIC] %v", r)
				err := c.Status(fiber.StatusInternalServerError).JSON(exc.InternalServerError("PANIC"))
				if err != nil {
					panicMiddlewareLogger.Error(err, "")
				}
				return
			}
		}()
		return c.Next()
	}
}
