package route

import (
	"github.com/haily-id/engine/internal/delivery/http/handler/company"
	"github.com/haily-id/engine/internal/delivery/http/handler/user"
	"github.com/haily-id/engine/internal/delivery/http/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type RouteConfig struct {
	UserHandler    *user.UserHandler
	CompanyHandler *company.CompanyHandler
	JWTSecret      string
}

func Setup(e *echo.Echo, cfg RouteConfig) {
	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())
	e.Use(echoMiddleware.RequestID())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "healthy",
		})
	})

	// API v1
	v1 := e.Group("/api/v1")

	// Public routes
	auth := v1.Group("/auth")
	auth.POST("/register", cfg.UserHandler.Register)
	auth.POST("/login", cfg.UserHandler.Login)

	// Protected user routes
	users := v1.Group("/users")
	users.Use(middleware.JWTAuth(cfg.JWTSecret))
	users.GET("/me", cfg.UserHandler.GetMe)
	users.GET("/me/companies", cfg.UserHandler.GetWithCompanies)
	users.PUT("/me", cfg.UserHandler.Update)
	users.DELETE("/me", cfg.UserHandler.Delete)
	users.GET("/:id", cfg.UserHandler.GetByID)
	users.GET("", cfg.UserHandler.List)
	users.POST("/companies/:company_id/join", cfg.UserHandler.JoinCompany)
	users.DELETE("/companies/:company_id/leave", cfg.UserHandler.LeaveCompany)

	// Protected company routes
	companies := v1.Group("/companies")
	companies.Use(middleware.JWTAuth(cfg.JWTSecret))
	companies.POST("", cfg.CompanyHandler.Create)
	companies.GET("/:id", cfg.CompanyHandler.GetByID)
	companies.GET("/code/:code", cfg.CompanyHandler.GetByCode)
	companies.PUT("/:id", cfg.CompanyHandler.Update)
	companies.DELETE("/:id", cfg.CompanyHandler.Delete)
	companies.GET("", cfg.CompanyHandler.List)
}
