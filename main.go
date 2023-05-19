package main

import (
	"fmt"

	"github.com/sandromai/go-http-server/models"
)

func main() {
	adminModel := models.Admin{}

	err := adminModel.Update(
		"New test",
		"test",
		"1234",
		2,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Admin #%v updated.\n", 2)

	adminId, err := adminModel.Authenticate("test", "1234")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Admin #%v authenticated.\n", adminId)

	admin, err := adminModel.FindById(adminId)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(*admin)
}
