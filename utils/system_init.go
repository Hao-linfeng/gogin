package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	RED *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")    //配置文件名 不带扩展格式
	viper.AddConfigPath("config") //路径
	err := viper.ReadInConfig()   //读取文件
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config mysql:", viper.Get("mysql"))
	fmt.Println("config mysql:", viper.Get("redis"))

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

func InitRedis() {
	RED = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.password"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

}

const (
	PublishKey = "websocket"
)

// 发布消息
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("publish .....", msg)
	err = RED.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// 订阅消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := RED.Subscribe(ctx, channel)
	fmt.Println("Subscribe........", ctx)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("PSubscribe .....", msg.Payload)
	return msg.Payload, err
}
