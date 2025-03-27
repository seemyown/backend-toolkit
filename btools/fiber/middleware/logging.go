package middleware

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

var requestLogger = log.NewSubLogger("request")

func LoggingMiddleware(c *fiber.Ctx) error {
	//Логирование входящих запросов
	tmStart := time.Now()
	err := c.Next()
	duration := time.Since(tmStart)

	requestLogger.Info(
		"method=%s path=%s ip=%s status=%d duration=%s queries=%v",
		c.Method(),
		c.Path(),
		c.IP(),
		c.Response().StatusCode(),
		duration,
		c.Queries(),
	)

	if err != nil {
		requestLogger.Error(err, "request error")
	}
	return err
}
