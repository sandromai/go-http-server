package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sandromai/go-http-server/models"
	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

func main() {
	currentPath, err := os.Executable()

	if err != nil {
		panic(err)
	}

	logger := &utils.Logger{
		FolderPath:     filepath.Join(filepath.Dir(currentPath), "logs"),
		FileName:       "main.log",
		MessagesPrefix: "[MAIN]:",
	}

	adminModel := &models.Admin{}

	admin, appErr := adminModel.Authenticate("john", "1234")

	if appErr != nil {
		appErr = logger.Save(
			fmt.Sprintf("(%v) %v\n", appErr.StatusCode, appErr.Message),
		)

		if appErr != nil {
			panic(appErr.Message)
		}

		return
	}

	jwt := &utils.JWT{}

	adminToken, appErr := jwt.Create(&types.AdminTokenPayload{
		AdminId:   admin.Id,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		CreatedAt: time.Now().Unix(),
	})

	if appErr != nil {
		appErr = logger.Save(
			fmt.Sprintf("(%v) %v\n", appErr.StatusCode, appErr.Message),
		)

		if appErr != nil {
			panic(appErr.Message)
		}

		return
	}

	fmt.Println(adminToken)

	authenticatedAdmin, appErr := utils.AuthenticateAdmin(adminToken)

	if appErr != nil {
		appErr = logger.Save(
			fmt.Sprintf("(%v) %v\n", appErr.StatusCode, appErr.Message),
		)

		if appErr != nil {
			panic(appErr.Message)
		}

		return
	}

	fmt.Println(*authenticatedAdmin)
}
