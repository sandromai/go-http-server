package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/sandromai/go-http-server/middlewares"
	"github.com/sandromai/go-http-server/models"
	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

type User struct {
	Timezone *time.Location
}

func (u *User) Authenticate(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "GET" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	user, token, appErr := middlewares.AuthenticateUser(
		request,
		u.Timezone,
	)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	utils.ReturnJSONResponse(
		writer,
		200,
		&struct {
			User  *types.User `json:"user"`
			Token string      `json:"token,omitempty"`
		}{User: user, Token: token},
	)
}

func (*User) Ban(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "PATCH" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	_, appErr := middlewares.AuthenticateAdmin(request)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	var userId string

	pathParts := strings.Split(request.URL.Path, "/")

	if pathParts[len(pathParts)-1] != "" {
		userId = pathParts[len(pathParts)-1]
	} else {
		userId = pathParts[len(pathParts)-2]
	}

	userModel := &models.User{}

	user, appErr := userModel.FindById(userId)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	if user.Banned {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "User is already banned.",
		})

		return
	}

	appErr = userModel.Ban(user.Id)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	utils.ReturnJSONResponse(
		writer,
		200,
		nil,
	)
}

func (*User) Unban(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "PATCH" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	_, appErr := middlewares.AuthenticateAdmin(request)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	var userId string

	pathParts := strings.Split(request.URL.Path, "/")

	if pathParts[len(pathParts)-1] != "" {
		userId = pathParts[len(pathParts)-1]
	} else {
		userId = pathParts[len(pathParts)-2]
	}

	userModel := &models.User{}

	user, appErr := userModel.FindById(userId)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	if !user.Banned {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "User is not banned.",
		})

		return
	}

	appErr = userModel.Unban(user.Id)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	utils.ReturnJSONResponse(
		writer,
		200,
		nil,
	)
}
