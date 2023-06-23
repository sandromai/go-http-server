package models

import (
	"database/sql"

	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

type User struct{}

func (*User) checkIdAvailability(
	id string,
) (
	bool,
	*types.AppError,
) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return false, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id` FROM `users` WHERE `id` = ? LIMIT 1",
	)

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to check ID availability.",
		}
	}

	defer statement.Close()

	adminId := ""

	err = statement.QueryRow(id).Scan(&adminId)

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

func (*User) checkEmailAvailability(
	email string,
) (
	bool,
	*types.AppError,
) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return false, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id` FROM `users` WHERE `email` = ? LIMIT 1",
	)

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to check email availability.",
		}
	}

	defer statement.Close()

	userId := ""

	err = statement.QueryRow(
		email,
	).Scan(
		&userId,
	)

	if err == sql.ErrNoRows {
		return true, nil
	}

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Error checking email availability.",
		}
	}

	return false, nil
}

func (user *User) generateId() (
	string,
	*types.AppError,
) {
	id, appErr := utils.GenerateUUIDv4()

	if appErr != nil {
		return "", appErr
	}

	idAvailability, appErr := user.checkIdAvailability(
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

		idAvailability, appErr = user.checkIdAvailability(
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

func (*User) FindById(
	id string,
) (
	*types.User,
	*types.AppError,
) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return nil, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id`, `email`, `banned`, `created_at` FROM `users` WHERE `id` = ? LIMIT 1",
	)

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to find user.",
		}
	}

	defer statement.Close()

	user := &types.User{}

	err = statement.QueryRow(id).Scan(
		&user.Id,
		&user.Email,
		&user.Banned,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, &types.AppError{
			StatusCode: 404,
			Message:    "User not found.",
		}
	}

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Error searching for user.",
		}
	}

	return user, nil
}

func (*User) FindByEmail(
	email string,
) (
	*types.User,
	*types.AppError,
) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return nil, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id`, `email`, `banned`, `created_at` FROM `users` WHERE `email` = ? LIMIT 1",
	)

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to find user.",
		}
	}

	defer statement.Close()

	user := &types.User{}

	err = statement.QueryRow(email).Scan(
		&user.Id,
		&user.Email,
		&user.Banned,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, &types.AppError{
			StatusCode: 404,
			Message:    "User not found.",
		}
	}

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Error searching for user.",
		}
	}

	return user, nil
}

func (user *User) Create(
	email string,
) (
	id string,
	appErr *types.AppError,
) {
	id, appErr = user.generateId()

	if appErr != nil {
		return "", appErr
	}

	emailIsAvailable, appErr := user.checkEmailAvailability(
		email,
	)

	if appErr != nil {
		return "", appErr
	}

	if !emailIsAvailable {
		return "", &types.AppError{
			StatusCode: 409,
			Message:    "Email already registered.",
		}
	}

	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return "", appErr
	}

	statement, err := dbConnection.Prepare(
		"INSERT INTO `users` (`id`, `email`) VALUES(?, ?)",
	)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to create user.",
		}
	}

	defer statement.Close()

	_, err = statement.Exec(
		id,
		email,
	)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Error creating user.",
		}
	}

	return id, nil
}

func (user *User) Ban(
	id string,
) *types.AppError {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return appErr
	}

	statement, err := dbConnection.Prepare("UPDATE `users` SET `banned` = 1 WHERE `id` = ?")

	if err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Failed to ban user.",
		}
	}

	defer statement.Close()

	if _, err = statement.Exec(id); err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Error banning user.",
		}
	}

	return nil
}

func (user *User) Unban(
	id string,
) *types.AppError {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return appErr
	}

	statement, err := dbConnection.Prepare("UPDATE `users` SET `banned` = 0 WHERE `id` = ?")

	if err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Failed to unban user.",
		}
	}

	defer statement.Close()

	if _, err = statement.Exec(id); err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Error unbanning user.",
		}
	}

	return nil
}
