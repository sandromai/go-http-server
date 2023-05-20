package main

import (
	"fmt"

	"github.com/sandromai/go-http-server/models"
)

func main() {
	adminModel := &models.Admin{}

	admin, err := adminModel.Authenticate("john", "123")

	if err != nil {
		fmt.Printf("(%v) %v\n", err.StatusCode, err.Message)
		return
	}

	fmt.Println(*admin)
}
