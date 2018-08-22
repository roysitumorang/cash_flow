package password

import (
	"cash_flow/service/user"
	"github.com/labstack/echo"
	"net/http"
)

func Reset(c echo.Context) error {
	var err error
	u, err := user.FindByEmail(c.FormValue("email"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not Found"})
	}
	err = u.ResetPassword()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Unprocessable Entity"})
	}
	return c.JSON(http.StatusOK, u)
}

func Save(c echo.Context) error {
	var err error
	u, err := user.FindByPasswordToken(c.Param("token"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not Found"})
	}
	u.Password = c.FormValue("password")
	u.PasswordConfirmation = c.FormValue("password_confirmation")
	if !u.Validate() {
		return c.JSON(http.StatusUnprocessableEntity, map[string]map[string]string{"errors": u.Errors})
	}
	err = u.SavePassword()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Unprocessable Entity"})
	}
	return c.JSON(http.StatusOK, u)
}
