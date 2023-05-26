package models

import (
	"database/sql"

	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

type LoginToken struct{}

func (*LoginToken) checkDeviceCodeAvailability(
	deviceCode string,
) (bool, *types.AppError) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return false, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id` FROM `login_tokens` WHERE `device_code` = ? LIMIT 1",
	)

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to check device code availability.",
		}
	}

	defer statement.Close()

	loginTokenId := int64(0)

	err = statement.QueryRow(
		deviceCode,
	).Scan(
		&loginTokenId,
	)

	if err == sql.ErrNoRows {
		return true, nil
	}

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Error checking device code availability.",
		}
	}

	return false, nil
}

func (loginToken *LoginToken) generateDeviceCode() (
	string,
	*types.AppError,
) {
	deviceCode, appErr := utils.GenerateUUIDv4()

	if appErr != nil {
		return "", appErr
	}

	deviceCodeAvailability, appErr := loginToken.checkDeviceCodeAvailability(
		deviceCode,
	)

	if appErr != nil {
		return "", appErr
	}

	for i := 0; i < 20 && !deviceCodeAvailability; i++ {
		deviceCode, appErr = utils.GenerateUUIDv4()

		if appErr != nil {
			return "", appErr
		}

		deviceCodeAvailability, appErr = loginToken.checkDeviceCodeAvailability(
			deviceCode,
		)

		if appErr != nil {
			return "", appErr
		}
	}

	if !deviceCodeAvailability {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to generate device code.",
		}
	}

	return deviceCode, nil
}
