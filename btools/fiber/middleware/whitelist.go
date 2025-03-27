package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/seemyown/backend-toolkit/btools/ext"
)

var ipLogger = log.NewSubLogger("ip_whitelist_middleware")

func WhilelistMiddleware(allowedIPs, allowedHosts []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		host := c.Hostname()
		ipLogger.Info("Incoming request from IP: %s; Host: %s", ip, host)
		if ext.Contains(allowedIPs, ip) || ext.Contains(allowedHosts, host) {
			ipLogger.Info("IP in whitelist")
			return c.Next()
		}
		ipLogger.Info("IP missing in whitelist")
		return c.SendStatus(fiber.StatusForbidden)
	}
}
