package models

import (
	"database/sql"

	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

type UserToken struct{}

func (*UserToken) checkLoginTokenAvailability(
	loginTokenId string,
) (bool, *types.AppError) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return false, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id` FROM `user_tokens` WHERE `from_login_token` = ? LIMIT 1",
	)

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to check login token availability.",
		}
	}

	defer statement.Close()

	userTokenId := ""

	err = statement.QueryRow(loginTokenId).Scan(&userTokenId)

	if err == sql.ErrNoRows {
		return true, nil
	}

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Error checking login token availability.",
		}
	}

	return false, nil
}

func (*UserToken) checkIdAvailability(
	id string,
) (bool, *types.AppError) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return false, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id` FROM `user_tokens` WHERE `id` = ? LIMIT 1",
	)

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to check ID availability.",
		}
	}

	defer statement.Close()

	userTokenId := ""

	err = statement.QueryRow(id).Scan(&userTokenId)

	if err == sql.ErrNoRows {
		return true, nil
	}

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Error checking ID availability.",
		}
	}

	return false, nil
}

func (userToken *UserToken) generateId() (
	string,
	*types.AppError,
) {
	id, appErr := utils.GenerateUUIDv4()

	if appErr != nil {
		return "", appErr
	}

	idAvailability, appErr := userToken.checkIdAvailability(
		id,
	)

	if appErr != nil {
		return "", appErr
	}

	for i := 0; i < 20 && !idAvailability; i++ {
		id, appErr = utils.GenerateUUIDv4()

		if appErr != nil {
			return "", appErr
		}

		idAvailability, appErr = userToken.checkIdAvailability(
			id,
		)

		if appErr != nil {
			return "", appErr
		}
	}

	if !idAvailability {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to generate ID.",
		}
	}

	return id, nil
}

func (*UserToken) FindById(
	id string,
) (
	*types.UserToken,
	*types.AppError,
) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return nil, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id`, `user_id`, `from_login_token`, `from_user_token`, `ip_address`, `device`, `disconnected`, `last_activity`, `expires_at`, `created_at` FROM `user_tokens` WHERE `id` = ? LIMIT 1",
	)

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to find user token.",
		}
	}

	defer statement.Close()

	userToken := &types.UserToken{}

	err = statement.QueryRow(id).Scan(
		&userToken.Id,
		&userToken.UserId,
		&userToken.FromLoginToken,
		&userToken.FromUserToken,
		&userToken.IPAddress,
		&userToken.Device,
		&userToken.Disconnected,
		&userToken.LastActivity,
		&userToken.ExpiresAt,
		&userToken.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, &types.AppError{
			StatusCode: 404,
			Message:    "User token not found.",
		}
	}

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Error searching user token.",
		}
	}

	return userToken, nil
}

func (userToken *UserToken) Create(
	userId string,
	fromLoginToken,
	fromUserToken *string,
	ipAddress,
	device string,
	expiresIn int64,
) (
	id string,
	appErr *types.AppError,
) {
	if fromLoginToken == nil && fromUserToken == nil {
		return "", &types.AppError{
			StatusCode: 400,
			Message:    "No login or user token provided.",
		}
	}

	if fromLoginToken != nil {
		loginTokenAvailable, appErr := userToken.checkLoginTokenAvailability(
			*fromLoginToken,
		)

		if appErr != nil {
			return "", appErr
		}

		if !loginTokenAvailable {
			return "", &types.AppError{
				StatusCode: 400,
				Message:    "Login token already used.",
			}
		}
	}

	id, appErr = userToken.generateId()

	if appErr != nil {
		return "", appErr
	}

	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return "", appErr
	}

	statement, err := dbConnection.Prepare(
		"INSERT INTO `user_tokens` (`id`, `user_id`, `from_login_token`, `from_user_token`, `ip_address`, `device`, `expires_at`) VALUES(?, ?, ?, ?, ?, ?, DATE_ADD(NOW(), INTERVAL ? SECOND))",
	)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to create user token.",
		}
	}

	defer statement.Close()

	_, err = statement.Exec(
		id,
		userId,
		fromLoginToken,
		fromUserToken,
		ipAddress,
		device,
		expiresIn,
	)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Error creating user token.",
		}
	}

	return id, nil
}

func (*UserToken) UpdateActivity(
	id string,
) *types.AppError {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return appErr
	}

	statement, err := dbConnection.Prepare(
		"UPDATE `user_tokens` SET `last_activity` = NOW() WHERE `id` = ?",
	)

	if err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Failed to update user token activity.",
		}
	}

	defer statement.Close()

	if _, err = statement.Exec(id); err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Error updating user token activity.",
		}
	}

	return nil
}

func (*UserToken) Disconnect(
	id string,
) *types.AppError {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return appErr
	}

	statement, err := dbConnection.Prepare(
		"UPDATE `user_tokens` SET `disconnected` = 1 WHERE `id` = ?",
	)

	if err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Failed to disconnect user token.",
		}
	}

	defer statement.Close()

	if _, err = statement.Exec(id); err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Error disconnecting user token.",
		}
	}

	return nil
}
