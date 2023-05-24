package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sandromai/go-http-server/models"
	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "POST" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	var body *loginData

	err := json.NewDecoder(request.Body).Decode(&body)

	if err == io.EOF {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "No data provided.",
		})

		return
	}

	if err != nil {
		utils.ReturnJSONResponse(writer, 500, nil)

		return
	}

	if body.Username == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert your username.",
		})

		return
	}

	if body.Password == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert your password.",
		})

		return
	}

	admin, appErr := (&models.Admin{}).Authenticate(
		body.Username,
		body.Password,
	)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	adminToken, appErr := (&utils.JWT{}).Create(&types.AdminTokenPayload{
		AdminId:   admin.Id,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		CreatedAt: time.Now().Unix(),
	})

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
			Admin *types.Admin `json:"admin"`
			Token string       `json:"token"`
		}{Admin: admin, Token: adminToken},
	)
}
