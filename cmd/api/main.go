package main

import (
	"github.com/synthe102/pgds-controller/internal/api/router"
)

func main() {
	r := router.New()

	err := r.Run()
	if err != nil {
		panic(err)
	}
}
