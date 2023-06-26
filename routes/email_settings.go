package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sandromai/go-http-server/middlewares"
	"github.com/sandromai/go-http-server/models"
	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

type EmailSetting struct{}

func (*EmailSetting) List(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "GET" {
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

	emailSettings, appErr := (&models.EmailSetting{}).List()

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
		emailSettings,
	)
}

func (*EmailSetting) Update(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "PUT" {
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

	var body *struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(request.Body).Decode(&body)

	if err == io.EOF {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert the host.",
		})

		return
	}

	if err != nil {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid data.",
		})

		return
	}

	if body.Host == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert the host.",
		})

		return
	}

	if body.Port == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert the port.",
		})

		return
	}

	if body.Username == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert the username.",
		})

		return
	}

	appErr = (&models.EmailSetting{}).Update(map[string]string{
		"host":     body.Host,
		"port":     body.Port,
		"username": body.Username,
		"password": body.Password,
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
		nil,
	)
}
