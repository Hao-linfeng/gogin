package main

import (
	"gogin/router"
	"gogin/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	r := router.Router()
	r.Run(":8000")
}
