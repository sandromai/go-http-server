package models

import (
	"database/sql"

	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

type LoginToken struct{}

func (*LoginToken) checkIdAvailability(
	id string,
) (bool, *types.AppError) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return false, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id` FROM `login_tokens` WHERE `id` = ? LIMIT 1",
	)

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to check ID availability.",
		}
	}

	defer statement.Close()

	loginTokenId := ""

	err = statement.QueryRow(id).Scan(&loginTokenId)

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

func (loginToken *LoginToken) generateId() (
	string,
	*types.AppError,
) {
	id, appErr := utils.GenerateUUIDv4()

	if appErr != nil {
		return "", appErr
	}

	idAvailability, appErr := loginToken.checkIdAvailability(
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

		idAvailability, appErr = loginToken.checkIdAvailability(
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

func (*LoginToken) FindById(
	id string,
) (*types.LoginToken, *types.AppError) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return nil, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id`, `email`, `ip_address`, `device`, `authorized`, `denied`, `expires_at`, `created_at` FROM `login_tokens` WHERE `id` = ? LIMIT 1",
	)

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to find login token.",
		}
	}

	defer statement.Close()

	loginToken := &types.LoginToken{}

	err = statement.QueryRow(id).Scan(
		&loginToken.Id,
		&loginToken.Email,
		&loginToken.IPAddress,
		&loginToken.Device,
		&loginToken.Authorized,
		&loginToken.Denied,
		&loginToken.ExpiresAt,
		&loginToken.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, &types.AppError{
			StatusCode: 404,
			Message:    "Login token not found.",
		}
	}

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Error searching login token.",
		}
	}

	return loginToken, nil
}

func (loginToken *LoginToken) Create(
	email,
	ipAddress,
	device string,
	expiresIn int64,
) (
	id string,
	appErr *types.AppError,
) {
	id, appErr = loginToken.generateId()

	if appErr != nil {
		return "", appErr
	}

	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return "", appErr
	}

	statement, err := dbConnection.Prepare(
		"INSERT INTO `login_tokens` (`id`, `email`, `ip_address`, `device`, `expires_at`) VALUES(?, ?, ?, ?, DATE_ADD(NOW(), INTERVAL ? SECOND))",
	)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to create login token.",
		}
	}

	defer statement.Close()

	_, err = statement.Exec(
		id,
		email,
		ipAddress,
		device,
		expiresIn,
	)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Error creating login token.",
		}
	}

	return id, nil
}

func (loginToken *LoginToken) Authorize(
	id string,
) *types.AppError {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return appErr
	}

	statement, err := dbConnection.Prepare("UPDATE `login_tokens` SET `authorized` = 1 WHERE `id` = ?")

	if err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Failed to authorize login token.",
		}
	}

	defer statement.Close()

	if _, err = statement.Exec(id); err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Error authorizing login token.",
		}
	}

	return nil
}

func (loginToken *LoginToken) Deny(
	id string,
) *types.AppError {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return appErr
	}

	statement, err := dbConnection.Prepare("UPDATE `login_tokens` SET `denied` = 1 WHERE `id` = ?")

	if err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Failed to deny login token.",
		}
	}

	defer statement.Close()

	if _, err = statement.Exec(id); err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Error denying login token.",
		}
	}

	return nil
}

func (*LoginToken) CountActiveByEmail(
	email string,
) (int64, *types.AppError) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return 0, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT COUNT(`id`) FROM `login_tokens` WHERE `email` = ? AND `expires_at` > NOW()",
	)

	if err != nil {
		return 0, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to count active login tokens.",
		}
	}

	defer statement.Close()

	activeLoginTokens := int64(0)

	err = statement.QueryRow(email).Scan(
		&activeLoginTokens,
	)

	if err != nil {
		return 0, &types.AppError{
			StatusCode: 500,
			Message:    "Error counting active login tokens.",
		}
	}

	return activeLoginTokens, nil
}

func (*LoginToken) GetLastCreationTimeByEmail(
	email string,
) (string, *types.AppError) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return "", appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `created_at` FROM `login_tokens` WHERE `email` = ? ORDER BY `id` DESC LIMIT 1",
	)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to get creation time from last login token.",
		}
	}

	defer statement.Close()

	creationTime := ""

	err = statement.QueryRow(email).Scan(
		&creationTime,
	)

	if err == sql.ErrNoRows {
		return "", nil
	}

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Error getting creation time from last login token.",
		}
	}

	return creationTime, nil
}
