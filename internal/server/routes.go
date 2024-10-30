package server

import (
	"net/http"

	"github.com/fastbyt3/diy-mssql-test/internal/handlers"
	"github.com/fastbyt3/diy-mssql-test/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			utils.Logger.Info(
				"request",
				zap.String("URI", v.URI),
				zap.Int("Status", v.Status),
			)
			return nil
		},
	}))

	e.Use(middleware.Recover())

	// Health endpoints
	healthHandlers := handlers.NewHealthHandler(s.db)
	e.GET("/health", healthHandlers.ApiHealth)
	e.GET("/db-health", healthHandlers.DbHealth)

	// User CRUD endpoints
	userHandlers := handlers.NewUserHandler(s.userService)
	e.GET("/all-users", userHandlers.GetAllUsers)
	e.GET("/user/:id", userHandlers.GetUser)
	e.GET("/user", userHandlers.CreateUser)

	return e
}
