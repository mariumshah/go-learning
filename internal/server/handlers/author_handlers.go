package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/playground/userapi/internal/repository"
)

type AuthorHandler struct {
	lib repository.LibraryRepo
}

func NewAuthorHandler(lib repository.LibraryRepo) *AuthorHandler {
	return &AuthorHandler{
		lib: lib,
	}
}

type addAuthorReq struct {
	Name string `json:"name"`
}

func (h *AuthorHandler) AddAuthor(c echo.Context) error {
	var req addAuthorReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}
	id, err := h.lib.AddAuthor(c.Request().Context(), req.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"id": id, "name": req.Name})
}
