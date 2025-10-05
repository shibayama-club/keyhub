package health

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/shibayama-club/keyhub/internal/domain/healthcheck"
)

type Handler struct {
	checkers []healthcheck.HealthChecker
}

func New(checkers ...healthcheck.HealthChecker) *Handler {
	return &Handler{
		checkers: checkers,
	}
}

func (h *Handler) Check(c echo.Context) error {
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		hasErr bool
	)

	for _, checker := range h.checkers {
		wg.Add(1)
		go func(c healthcheck.HealthChecker) {
			defer wg.Done()
			if err := c.Ping(context.Background()); err != nil {
				mu.Lock()
				defer mu.Unlock()
				hasErr = true
				slog.Error("health check failed",
					"checker", c.Name(),
					"error", err,
				)
			}
		}(checker)
	}

	wg.Wait()

	if hasErr {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "not_serving",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "serving",
	})
}
