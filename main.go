package main

import (
	"github.com/seeu359/go-project-278/router"
)

func main() {
	r := router.New()

	r.Run(":8080")
}
