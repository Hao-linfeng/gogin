package models

import (
	"fmt"
	"gogin/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string
	Email         string
	Identidy      string
	ClientIp      string
	ClientPort    string
	LoginTime     time.Time
	HeartBeatTime time.Time
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"logon_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func CreateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user *UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, Password: user.Password})
}
