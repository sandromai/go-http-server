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

type UserToken struct {
	Timezone *time.Location
}

func (u *UserToken) Disconnect(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "PATCH" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	user, _, appErr := middlewares.AuthenticateUser(
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

	var userTokenId string

	pathParts := strings.Split(request.URL.Path, "/")

	if pathParts[len(pathParts)-1] != "" {
		userTokenId = pathParts[len(pathParts)-1]
	} else {
		userTokenId = pathParts[len(pathParts)-2]
	}

	userTokenModel := &models.UserToken{}

	userToken, appErr := userTokenModel.FindById(
		userTokenId,
	)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	if user.Id != userToken.UserId {
		utils.ReturnJSONResponse(writer, 403, &types.ReturnError{
			Error: "Unauthorized action.",
		})

		return
	}

	if userToken.Disconnected {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "This token was already disconnected.",
		})

		return
	}

	appErr = userTokenModel.Disconnect(userToken.Id)

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
