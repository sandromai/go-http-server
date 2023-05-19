package models

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/sandromai/go-http-server/types"
)

type Admin struct{}

func (Admin) checkUsernameAvailability(
	username string,
	excludeId int64,
) (bool, error) {
	dbConnection, err := getDBConnection()

	if err != nil {
		return false, err
	}

	defer dbConnection.Close()

	statement, err := dbConnection.Prepare(
		"SELECT `id` FROM `admins` WHERE `username` = ? AND `id` != ? LIMIT 1",
	)

	if err != nil {
		return false, err
	}

	defer statement.Close()

	var adminId int64

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
		return false, err
	}

	return false, nil
}

func (Admin) FindById(
	id int64,
) (*types.Admin, error) {
	dbConnection, err := getDBConnection()

	if err != nil {
		return nil, err
	}

	defer dbConnection.Close()

	statement, err := dbConnection.Prepare(
		"SELECT `id`, `name`, `username`, `created_at` FROM `admins` WHERE `id` = ? LIMIT 1",
	)

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	var admin types.Admin

	err = statement.QueryRow(id).Scan(
		&admin.Id,
		&admin.Name,
		&admin.Username,
		&admin.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (admin Admin) Create(
	name,
	username,
	password string,
) (int64, error) {
	usernameIsAvailable, err := admin.checkUsernameAvailability(
		username,
		0,
	)

	if err != nil {
		return 0, err
	}

	if !usernameIsAvailable {
		return 0, errors.New("username already registered")
	}

	dbConnection, err := getDBConnection()

	if err != nil {
		return 0, err
	}

	defer dbConnection.Close()

	statement, err := dbConnection.Prepare(
		"INSERT INTO `admins` (`name`, `username`, `password`) VALUES(?, ?, ?)",
	)

	if err != nil {
		return 0, err
	}

	defer statement.Close()

	passwordBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		12,
	)

	if err != nil {
		return 0, err
	}

	encryptedPassword := string(passwordBytes)

	result, err := statement.Exec(
		name,
		username,
		encryptedPassword,
	)

	if err != nil {
		return 0, err
	}

	adminId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return adminId, nil
}

func (admin Admin) Update(
	name,
	username,
	password string,
	id int64,
) error {
	usernameIsAvailable, err := admin.checkUsernameAvailability(
		username,
		id,
	)

	if err != nil {
		return err
	}

	if !usernameIsAvailable {
		return errors.New("username already registered")
	}

	dbConnection, err := getDBConnection()

	if err != nil {
		return err
	}

	defer dbConnection.Close()

	query := "UPDATE `admins` SET `name` = ?, `username` = ?"

	if password != "" {
		query += ", `password` = ?"
	}

	query += " WHERE `id` = ?"

	statement, err := dbConnection.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	var result sql.Result

	if password != "" {
		passwordBytes, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			12,
		)

		if err != nil {
			return err
		}

		encryptedPassword := string(passwordBytes)

		result, err = statement.Exec(
			name,
			username,
			encryptedPassword,
			id,
		)
	} else {
		result, err = statement.Exec(
			name,
			username,
			id,
		)
	}

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("admin not found")
	}

	return nil
}

func (Admin) Authenticate(
	username,
	password string,
) (int64, error) {
	dbConnection, err := getDBConnection()

	if err != nil {
		return 0, err
	}

	defer dbConnection.Close()

	statement, err := dbConnection.Prepare(
		"SELECT `id`, `password` FROM `admins` WHERE `username` = ? LIMIT 1",
	)

	if err != nil {
		return 0, err
	}

	defer statement.Close()

	var adminId int64
	var adminPassword string

	err = statement.QueryRow(username).Scan(
		&adminId,
		&adminPassword,
	)

	if err == sql.ErrNoRows {
		return 0, errors.New("incorrect username or password")
	}

	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(adminPassword),
		[]byte(password),
	)

	if err != nil {
		return 0, errors.New("incorrect username or password")
	}

	return adminId, nil
}
