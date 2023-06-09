package models

import (
	"fmt"
	"gogin/utils"
	"time"

	"gorm.io/gorm"
)

// 用户
type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9})"`
	Email         string `valid:"email"`
	Identidy      string
	ClientIp      string
	ClientPort    string
	LoginTime     time.Time
	HeartBeatTime time.Time
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"logon_out_time"`
	IsLogout      bool
	DeviceInfo    string
	Salt          string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func FindUserByNameAndPwd(name string, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and password=?", name, password).First(&user)
	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identidy", temp)
	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
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
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, Password: user.Password, Phone: user.Phone, Email: user.Email})
}
