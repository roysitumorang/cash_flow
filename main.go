package main

import (
	"cash_flow/controller/account"
	"cash_flow/controller/password"
	"cash_flow/controller/users"
	"cash_flow/util/conn"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := conn.InitDB()
	handleError(err)
	defer conn.DB.Close()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	u := e.Group("/api/v1")
	u.GET("/users", users.Index)
	u.POST("/users", users.Create)
	u.GET("/users/:id", users.Show)
	u.PUT("/users/:id", users.Update)
	u.PUT("/users/:token/activate", users.Activate)
	u.DELETE("/users/:id", users.Destroy)
	u.POST("/account/create", account.Create)
	u.PUT("/account/:token/activate", account.Activate)
	u.PUT("/account/login", account.Login)
	u.PUT("/password/reset", password.Reset)
	u.PUT("/password/:token/save", password.Save)
	r := e.Group("/api/v1/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("/account", account.Show)
	r.PUT("/account", account.Update)
	e.Logger.Fatal(e.Start(":3000"))
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
