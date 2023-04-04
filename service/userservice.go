package service

import (
	"gogin/models"

	"github.com/gin-gonic/gin"
)

// GetUserList
// @Tags 首页
// @Accept json
// @Produce json
// @Success 200 {string} json{"code", "message"}
// @Router /user/GetUserList [get]
func GetUserList(c *gin.Context) {
	userList := models.GetUserList()
	c.JSON(200, userList)
}
