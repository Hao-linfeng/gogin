package main

import (
	"gogin/router"
)

func main() {
	r := router.Router()
	r.Run(":8000")
}
