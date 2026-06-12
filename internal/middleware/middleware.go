package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestID injects a unique X-Request-ID header into responses and context if not present.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		rid := c.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Set("X-Request-ID", rid)
		c.Locals("requestId", rid)
		return c.Next()
	}
}

// RequestLogger logs HTTP requests including their method, path, response status, and duration.
func RequestLogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Execute request handler chain
		err := c.Next()

		duration := time.Since(start)
		rid := c.Locals("requestId")

		fields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
		}

		if rid != nil {
			fields = append(fields, zap.Any("requestId", rid))
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Error("Request failed", fields...)
		} else {
			logger.Info("Request completed", fields...)
		}

		return err
	}
}
