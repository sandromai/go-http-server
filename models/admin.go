package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

type Admin struct{}

func (*Admin) checkIdAvailability(
	id string,
) (bool, *types.AppError) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return false, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id` FROM `admins` WHERE `id` = ? LIMIT 1",
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

func (*Admin) checkUsernameAvailability(
	username string,
	excludeId string,
) (bool, *types.AppError) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return false, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id` FROM `admins` WHERE `username` = ? AND `id` != ? LIMIT 1",
	)

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to check username availability.",
		}
	}

	defer statement.Close()

	adminId := ""

	err = statement.QueryRow(
		username,
		excludeId,
	).Scan(
		&adminId,
	)

	if err == sql.ErrNoRows {
		return true, nil
	}

	if err != nil {
		return false, &types.AppError{
			StatusCode: 500,
			Message:    "Error checking username availability.",
		}
	}

	return false, nil
}

func (admin *Admin) generateId() (
	string,
	*types.AppError,
) {
	id, appErr := utils.GenerateUUIDv4()

	if appErr != nil {
		return "", appErr
	}

	idAvailability, appErr := admin.checkIdAvailability(
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

		idAvailability, appErr = admin.checkIdAvailability(
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

func (*Admin) FindById(
	id string,
) (*types.Admin, *types.AppError) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return nil, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id`, `name`, `username`, `created_by`, `created_at` FROM `admins` WHERE `id` = ? LIMIT 1",
	)

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to find admin.",
		}
	}

	defer statement.Close()

	admin := &types.Admin{}

	err = statement.QueryRow(id).Scan(
		&admin.Id,
		&admin.Name,
		&admin.Username,
		&admin.CreatedBy,
		&admin.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, &types.AppError{
			StatusCode: 404,
			Message:    "Admin not found.",
		}
	}

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Error searching for admin.",
		}
	}

	return admin, nil
}

func (admin *Admin) Create(
	name,
	username,
	password string,
	createdBy *string,
) (
	id string,
	appErr *types.AppError,
) {
	id, appErr = admin.generateId()

	if appErr != nil {
		return "", appErr
	}

	usernameIsAvailable, appErr := admin.checkUsernameAvailability(
		username,
		"",
	)

	if appErr != nil {
		return "", appErr
	}

	if !usernameIsAvailable {
		return "", &types.AppError{
			StatusCode: 409,
			Message:    "Username already registered.",
		}
	}

	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return "", appErr
	}

	statement, err := dbConnection.Prepare(
		"INSERT INTO `admins` (`id`, `name`, `username`, `password`, `created_by`) VALUES(?, ?, ?, ?, ?)",
	)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to create admin.",
		}
	}

	defer statement.Close()

	passwordBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		12,
	)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to hash password.",
		}
	}

	encryptedPassword := string(passwordBytes)

	_, err = statement.Exec(
		id,
		name,
		username,
		encryptedPassword,
		createdBy,
	)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Error creating admin.",
		}
	}

	return id, nil
}

func (admin *Admin) Update(
	name,
	username,
	password,
	id string,
) *types.AppError {
	usernameIsAvailable, appErr := admin.checkUsernameAvailability(
		username,
		id,
	)

	if appErr != nil {
		return appErr
	}

	if !usernameIsAvailable {
		return &types.AppError{
			StatusCode: 409,
			Message:    "Username already registered.",
		}
	}

	query := "UPDATE `admins` SET `name` = ?, `username` = ?"

	if password != "" {
		query += ", `password` = ?"
	}

	query += " WHERE `id` = ?"

	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return appErr
	}

	statement, err := dbConnection.Prepare(query)

	if err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Failed to update admin.",
		}
	}

	defer statement.Close()

	if password != "" {
		passwordBytes, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			12,
		)

		if err != nil {
			return &types.AppError{
				StatusCode: 500,
				Message:    "Failed to hash password.",
			}
		}

		encryptedPassword := string(passwordBytes)

		_, err = statement.Exec(
			name,
			username,
			encryptedPassword,
			id,
		)

		if err != nil {
			return &types.AppError{
				StatusCode: 500,
				Message:    "Error updating admin.",
			}
		}
	} else {
		_, err = statement.Exec(
			name,
			username,
			id,
		)

		if err != nil {
			return &types.AppError{
				StatusCode: 500,
				Message:    "Error updating admin.",
			}
		}
	}

	return nil
}

func (*Admin) Authenticate(
	username,
	password string,
) (*types.Admin, *types.AppError) {
	dbConnection, appErr := getDBInstance()

	if appErr != nil {
		return nil, appErr
	}

	statement, err := dbConnection.Prepare(
		"SELECT `id`, `name`, `username`, `password`, `created_by`, `created_at` FROM `admins` WHERE `username` = ? LIMIT 1",
	)

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to authenticate admin.",
		}
	}

	defer statement.Close()

	admin := &types.Admin{}
	adminPassword := ""

	err = statement.QueryRow(username).Scan(
		&admin.Id,
		&admin.Name,
		&admin.Username,
		&adminPassword,
		&admin.CreatedBy,
		&admin.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, &types.AppError{
			StatusCode: 401,
			Message:    "Incorrect username or password.",
		}
	}

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Error authenticating admin.",
		}
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(adminPassword),
		[]byte(password),
	)

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 401,
			Message:    "Incorrect username or password.",
		}
	}

	return admin, nil
}
