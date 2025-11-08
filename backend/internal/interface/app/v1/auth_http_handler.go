package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GoogleLogin(c echo.Context) error {
	ctx := c.Request().Context()

	authURL, err := h.useCase.StartGoogleLogin(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start login process")
	}

	return c.Redirect(http.StatusFound, authURL)
}

func (h *Handler) GoogleCallback(c echo.Context) error {
	ctx := c.Request().Context()

	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" || state == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid callback parameters")
	}

	sessionID, err := h.useCase.GoogleCallback(ctx, code, state)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication failed")
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60,
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/app")
}
