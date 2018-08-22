package users

import (
	"cash_flow/service/user"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func Index(c echo.Context) error {
	var (
		page int
		err  error
	)
	term := c.QueryParam("term")
	page, err = strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	users, err := user.Where(term, page)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

func Create(c echo.Context) error {
	u := user.New()
	u.Name = c.FormValue("name")
	u.Email = c.FormValue("email")
	u.Password = c.FormValue("password")
	if err := u.Create(); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
}

func Show(c echo.Context) error {
	var err error
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not Found"})
	}
	u, err := user.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not Found"})
	}
	return c.JSON(http.StatusOK, u)
}

func Update(c echo.Context) error {
	var err error
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not Found"})
	}
	u, err := user.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not Found"})
	}
	u.Name = c.FormValue("name")
	u.Email = c.FormValue("email")
	u.Password = c.FormValue("password")
	err = u.Update()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Unprocessable Entity"})
	}
	return c.JSON(http.StatusOK, u)
}

func Destroy(c echo.Context) error {
	var err error
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not Found"})
	}
	u, err := user.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not Found"})
	}
	err = u.Destroy()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Unprocessable Entity"})
	}
	return c.NoContent(http.StatusNoContent)
}