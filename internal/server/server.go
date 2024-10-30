package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fastbyt3/diy-mssql-test/pkg/database"
	"github.com/fastbyt3/diy-mssql-test/pkg/repository"
	"github.com/fastbyt3/diy-mssql-test/pkg/service"
	"github.com/fastbyt3/diy-mssql-test/pkg/utils"
)

type Server struct {
	baseRoute   string
	port        int
	db          database.Service
	userService *service.UserService
}

func New() *http.Server {
	port, err := strconv.Atoi(utils.GetEnvVarOrDefault("APP_PORT", "9001"))
	if err != nil {
		utils.Logger.Fatal("APP_PORT is not a valid number")
	}

	dbService := database.New()
	userService := service.NewUserService(*repository.NewUserRepo(dbService.Db))

	s := &Server{
		baseRoute:   utils.GetEnvVarOrDefault("BASE_ROUTE", "/"),
		port:        port,
		db:          dbService,
		userService: userService,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
