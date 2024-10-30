package handlers

import (
	"net/http"
	"strconv"

	"github.com/fastbyt3/diy-mssql-test/pkg/service"
	"github.com/fastbyt3/diy-mssql-test/pkg/utils"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	allUsers, err := h.userService.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewServerError("Unable to fetch all users", err.Error()))
	}

	return c.JSON(http.StatusOK, allUsers)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewServerError("invalid user id", "user id needs to be an integer"))
	}

	user, err := h.userService.GetUser(c.Request().Context(), int32(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewServerError("Failed to get user", err.Error()))
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")

	err := h.userService.CreateUser(c.Request().Context(), username, email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewServerError("failed to create user", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "true"})
}
