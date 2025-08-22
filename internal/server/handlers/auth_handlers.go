package handlers

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/labstack/echo/v4"
	"github.com/playground/userapi/internal/app/domain"
	"github.com/playground/userapi/internal/repository"
	"github.com/playground/userapi/internal/server/middleware"
	"github.com/playground/userapi/pkg/utils"
)

type AuthHandler struct {
	userRepo repository.UserRepo
	libRepo  repository.LibraryRepo
}

func NewAuthHandler(ur repository.UserRepo, lr repository.LibraryRepo) *AuthHandler {
	return &AuthHandler{userRepo: ur, libRepo: lr}
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req registerReq) Validate() error {
	return validation.ValidateStruct(
		&req,
		validation.Field(&req.Email, is.Email,
			validation.Required),
		validation.Field(&req.Password, validation.Min(8),
			validation.Required),
	)
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req registerReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "invalid body")
	}
	// hash password
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "hash error")
	}
	u := &domain.User{
		Email: req.Email,
		Hash:  hash,
	}
	// create user
	if err := h.userRepo.Create(context.Background(), u); err != nil {
		// detect email exists
		if err == repository.ErrEmailExists {
			return echo.NewHTTPError(http.StatusConflict, "email already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// create library for user (best-effort)
	// _ = h.libRepo.CreateLibraryForUser(context.Background(), u.ID)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":    u.ID,
		"email": u.Email,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req loginReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			"invalid body")
	}
	dbu, err := h.userRepo.GetByEmail(context.Background(), req.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}
	if err := utils.CheckPassword(dbu.Hash, req.Password); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}
	token, err := utils.GenerateToken(dbu.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "token error")
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) Verify(c echo.Context) error {

	uid := middleware.GetUserID(c)
	if uid == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": uid,
	})
}
