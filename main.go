package main

import (
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
	e.GET("/api/v1/users", users.Index)
	e.POST("/api/v1/users", users.Create)
	e.GET("/api/v1/users/:id", users.Show)
	e.PUT("/api/v1/users/:id", users.Update)
	e.DELETE("/api/v1/users/:id", users.Destroy)
	e.Logger.Fatal(e.Start(":3000"))
}

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
