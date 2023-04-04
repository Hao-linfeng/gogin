package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitConfig() {
	viper.SetConfigName("app")    //配置文件名 不带扩展格式
	viper.AddConfigPath("config") //路径
	err := viper.ReadInConfig()   //读取文件
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config mysql:", viper.Get("mysql"))

}
func InitMysql() {

	//自定义日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.Lshortfile),
		logger.Config{
			SlowThreshold: time.Second, //慢sql阈值
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	// user := models.UserBasic{}
	// DB.Find(user)
	// fmt.Println(user)
}
