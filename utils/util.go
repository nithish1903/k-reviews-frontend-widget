package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func Respond(c *gin.Context, status int, u interface{}) error {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(status)
	return json.NewEncoder(c.Writer).Encode(u)
}
