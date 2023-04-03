package service

import (
	"gogin/models"

	"github.com/gin-gonic/gin"
)

func GetUserList(c *gin.Context) {
	userList := models.GetUserList()
	c.JSON(200, userList)
}
