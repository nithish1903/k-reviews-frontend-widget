package utils

import (
	"encoding/json"

	"github.com/labstack/echo"
)

func Respond(c echo.Context, status int, u interface{}) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(status)
	return json.NewEncoder(c.Response()).Encode(u)
}
