package ui

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Register(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		g := e.Router.Group("/app")
		g.GET("", func(c echo.Context) error {
			return c.Redirect(http.StatusFound, "/app/feed")
		})
		g.GET("/login", handleLogin(app))
		g.POST("/login", handleLogin(app))
		g.GET("/signup", handleSignup(app))
		g.POST("/signup", handleSignup(app))
		g.GET("/logout", handleLogout)
		g.GET("/feed", handleFeed(app))
		g.GET("/listings/new", handleNewListing(app))
		g.POST("/listings/new", handleNewListing(app))
		return nil
	})
}
