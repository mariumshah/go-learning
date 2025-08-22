package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"github.com/playground/userapi/internal/config"
	"github.com/playground/userapi/internal/repository"
	rhandlers "github.com/playground/userapi/internal/server/handlers"
	"github.com/playground/userapi/internal/server/middleware"
	"github.com/playground/userapi/pkg/utils"
)

func main() {
	conf := config.Load()

	db, err := sql.Open("mysql", conf.DBURL)
	if err != nil {
		fmt.Printf("db open: %v", err)
		panic(fmt.Errorf("db open: %w", err))
	}
	defer db.Close()

	//create repos
	userRepo := repository.NewUserRepo(db)
	libRepo := repository.NewLibraryRepo(db)

	//create  handlers
	authH := rhandlers.NewAuthHandler(userRepo, libRepo)
	libH := rhandlers.NewLibraryHandler(libRepo)
	authorH := rhandlers.NewAuthorHandler(libRepo)

	e := echo.New()

	// public endpoints
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	e.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"version": utils.GetVersion()})
	})

	// e.POST("/authors", authorH.AddAuthor)

	e.POST("/register", authH.Register)
	e.POST("/login", authH.Login)

	// protected endpoints
	meGroup := e.Group("/auth")
	meGroup.Use(middleware.RequireAuth)
	meGroup.GET("/me", authH.Verify)

	libGroup := e.Group("/library")
	libGroup.Use(middleware.RequireAuth)
	libGroup.POST("/books", libH.AddBook)
	libGroup.GET("/books", libH.ListBooks) //list books in the user's library
	libGroup.GET("/all_books", libH.ListAllBooks)

	aGroup := e.Group("/authors")
	aGroup.Use(middleware.RequireAuth)
	aGroup.POST("", authorH.AddAuthor)

	e.Logger.Infof("Server starting on :%s", conf.Port)
	e.Logger.Fatal(e.Start(":" + conf.Port))

}
