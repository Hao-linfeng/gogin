package service

import (
	"fmt"
	"gogin/models"
	"gogin/utils"
	"math/rand"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Accept json
// @Produce json
// @Success 200 {string} json{"code", "message"}
// @Router /user/GetUserList [get]
func GetUserList(c *gin.Context) {
	userList := models.GetUserList()
	c.JSON(200, userList)
}

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Accept json
// @Produce json
// @Success 200 {string} json{"code", "message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	name := c.Query("name")
	password := c.Query("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "该用户不存在",
			"data":    data,
		})
		return

	}
	flag := utils.ValidPasswoed(password, user.Password, user.Salt)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "密码或者用户名不正确",
			"data":    data,
		})
		return
	}

	pwd := utils.MakePasswoed(password, user.Salt)
	models.FindUserByNameAndPwd(name, pwd)

	c.JSON(200, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    data,
	})
}

// GetUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/CreateUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "两次密码不一致",
		})
		return
	}
	n := models.FindUserByName(user.Name)
	if n.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户名已经注册",
		})
		return
	}
	e := models.FindUserByEmail(user.Email)

	if e.Email != "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "邮箱已经注册",
		})
		return
	}
	p := models.FindUserByPhone(user.Phone)
	if p.Phone != "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "手机号已经注册",
		})
		return
	}

	user.Password = utils.MakePasswoed(password, salt)
	models.CreateUser(&user)
	c.JSON(200, gin.H{
		"code":    0,
		"messahe": "新增用户成功",
		"data":    user,
	})

}

// GetUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "ID"
// @Success 200 {string} json{"code", "message"}
// @Router /user/DeleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(&user)
	c.JSON(200, gin.H{
		"code":    0,
		"messahe": "删除用户成功",
		"data":    user,
	})

}

// GetUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "ID"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param phone formData string false "手机号"
// @param email formData string false "邮件"
// @Success 200 {string} json{"code", "message"}
// @Router /user/UpdateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("Password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1,
			"messahe": "修改参数有误",
			"data":    user,
		})
	} else {
		models.UpdateUser(&user)
		c.JSON(200, gin.H{
			"code":    0,
			"messahe": "修改用户成功",
			"data":    user,
		})
	}

}
