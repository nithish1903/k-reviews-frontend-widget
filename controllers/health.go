package controllers

import (
	"k-reviews-frontend-api/entity"
	"k-reviews-frontend-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Health(c *gin.Context) {
	u := &entity.Response{
		Status:  200,
		Message: "Service is up and running.",
		Data:    make([]interface{}, 0),
		Error:   entity.Empty{},
	}
	zap.L().Info("Health request triggered!")
	utils.Respond(c, http.StatusOK, u)
}
