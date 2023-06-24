package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

type EmailSetting struct{}

func (*EmailSetting) List() (
	*types.EmailSetting,
	*types.AppError,
) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return nil, appErr
	}

	emailSettings := &types.EmailSetting{}

	err := dbConnection.QueryRow("SELECT `host`, `port`, `username`, `password` FROM `email_settings` ORDER BY `id` DESC LIMIT 1").Scan(
		&emailSettings.Host,
		&emailSettings.Port,
		&emailSettings.Username,
		&emailSettings.Password,
	)

	if err == sql.ErrNoRows {
		return nil, &types.AppError{
			StatusCode: 404,
			Message:    "Email settings not found.",
		}
	}

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Error searching for email settings.",
		}
	}

	emailSettings.Password, appErr = utils.Decrypt(emailSettings.Password)

	if appErr != nil {
		return nil, appErr
	}

	return emailSettings, nil
}

func (*EmailSetting) Update(
	data map[string]string,
) *types.AppError {
	var updates []string
	var values []any

	for column, value := range data {
		if column == "password" {
			if value == "" {
				continue
			}

			encryptedPassword, err := utils.Encrypt(value)

			if err != nil {
				return &types.AppError{
					StatusCode: 500,
					Message:    "Failed to encrypt password.",
				}
			}

			values = append(values, encryptedPassword)
		} else {
			values = append(values, value)
		}

		updates = append(updates, "`"+column+"` = ?")
	}

	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return appErr
	}

	statement, err := dbConnection.Prepare("UPDATE `email_settings` SET " + strings.Join(updates, ", "))

	if err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Failed to update email settings.",
		}
	}

	defer statement.Close()

	_, err = statement.Exec(values...)

	if err != nil {
		fmt.Println(err.Error())

		return &types.AppError{
			StatusCode: 500,
			Message:    "Error updating email settings.",
		}
	}

	return nil
}
