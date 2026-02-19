package route

import (
	authHandler "github.com/haily-id/engine/internal/delivery/http/handler/auth"
	"github.com/haily-id/engine/internal/delivery/http/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type RouteConfig struct {
	AuthHandler *authHandler.Handler
	JWTSecret   string
}

func Setup(e *echo.Echo, cfg RouteConfig) {
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())
	e.Use(echoMiddleware.RequestID())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy"})
	})

	v1 := e.Group("/api/v1")

	// ── Auth (public) ────────────────────────────────────────────
	auth := v1.Group("/auth")
	auth.POST("/register", cfg.AuthHandler.Register)
	auth.POST("/login", cfg.AuthHandler.Login)
	auth.POST("/verify-email", cfg.AuthHandler.VerifyEmail)
	auth.POST("/resend-otp", cfg.AuthHandler.ResendOTP)

	// ── Auth (protected) ─────────────────────────────────────────
	authProtected := v1.Group("/auth")
	authProtected.Use(middleware.JWTAuth(cfg.JWTSecret))
	authProtected.GET("/me", cfg.AuthHandler.GetMe)
}

