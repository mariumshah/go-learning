// package handlers

// import (
// 	"context"
// 	// "database/sql"
// 	"net/http"
// 	// "strconv"

// 	"github.com/labstack/echo/v4"
// 	"github.com/playground/userapi/internal/app/domain"
// 	"github.com/playground/userapi/internal/repository"
// 	"github.com/playground/userapi/internal/server/middleware"
// )

// type LibraryHandler struct {
// 	libRepo repository.LibraryRepo
// }

// func NewLibraryHandler(lr repository.LibraryRepo) *LibraryHandler {
// 	return &LibraryHandler{libRepo: lr}
// }

// type addBookReq struct {
// 	Title           string `json:"title"`
// 	Author          int `json:"author"`
// 	ISBN            string `json:"isbn,omitempty"`
// 	PublicationYear int    `json:"publication_year,omitempty"`
// }

// func (h *LibraryHandler) AddBook(c echo.Context) error {
// 	uid := middleware.GetUserID(c)
// 	if uid == 0 {
// 		return echo.NewHTTPError(http.StatusUnauthorized, "missing user")
// 	}
// 	var req addBookReq
// 	if err := c.Bind(&req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
// 	}
// 	b := &domain.Book{
// 		Title:           req.Title,
// 		Author:          req.Author,
// 		// ISBN:            req.ISBN,
// 		// AddedAt:         /* leave zero, converter will set on insert */,
// 	}
// 	created, err := h.libRepo.AddBook(context.Background(), uid, b)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 	}
// 	return c.JSON(http.StatusCreated, created)
// }

// func (h *LibraryHandler) ListBooks(c echo.Context) error {
// 	uid := middleware.GetUserID(c)
// 	if uid == 0 {
// 		return echo.NewHTTPError(http.StatusUnauthorized, "missing user")
// 	}
// 	books, err := h.libRepo.ListBooks(context.Background(), uid)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 	}
// 	return c.JSON(http.StatusOK, books)
// }

// // func (h *LibraryHandler) RemoveBook(c echo.Context) error {
// // 	uid := middleware.GetUserID(c)
// // 	if uid == 0 {
// // 		return echo.NewHTTPError(http.StatusUnauthorized, "missing user")
// // 	}
// // 	idS := c.Param("id")
// // 	id, err := strconv.Atoi(idS)
// // 	if err != nil {
// // 		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
// // 	}
// // 	if err := h.libRepo.RemoveBook(context.Background(), uid, id); err != nil {
// // 		if err == sql.ErrNoRows {
// // 			return echo.NewHTTPError(http.StatusNotFound, "book not found")
// // 		}
// // 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// // 	}
// // 	return c.NoContent(http.StatusNoContent)
// // }

package handlers

import (
	// "context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/playground/userapi/internal/app/domain"
	"github.com/playground/userapi/internal/repository"
	"github.com/playground/userapi/internal/server/middleware"
)

type LibraryHandler struct {
	libRepo repository.LibraryRepo
}

func NewLibraryHandler(lr repository.LibraryRepo) *LibraryHandler {
	return &LibraryHandler{libRepo: lr}
}

type addBookReq struct {
	Title           string `json:"title"`
	Author          int    `json:"author"`
	ISBN            string `json:"isbn,omitempty"`
	PublicationYear int    `json:"publication_year,omitempty"`
}

func (h *LibraryHandler) AddBook(c echo.Context) error {
	uid := middleware.GetUserID(c)
	if uid == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing user")
	}

	var req addBookReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	// basic validation
	if req.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "title is required")
	}
	if req.Author == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "author id is required")
	}

	// use the request context so cancellations/timeouts propagate
	ctx := c.Request().Context()

	b := &domain.Book{
		Title:  req.Title,
		Author: req.Author,
		// ISBN:            req.ISBN,
		// PublicationYear: req.PublicationYear,
	}

	created, err := h.libRepo.AddBook(ctx, uid, b)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, created)
}

func (h *LibraryHandler) ListBooks(c echo.Context) error {
	uid := middleware.GetUserID(c)
	if uid == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing user")
	}
	ctx := c.Request().Context()
	books, err := h.libRepo.ListBooks(ctx, uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}

func (h *LibraryHandler) ListAllBooks(c echo.Context) error {
	uid := middleware.GetUserID(c)
	if uid == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing user")
	}
	ctx := c.Request().Context()
	books, err := h.libRepo.ListAllBooks(ctx, uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, books)
}

func (h *LibraryHandler) RemoveBook(c echo.Context) error {
	uid := middleware.GetUserID(c)
	if uid == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing user")
	}
	idS := c.Param("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	ctx := c.Request().Context()
	if err := h.libRepo.RemoveBook(ctx, uid, id); err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "book not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
