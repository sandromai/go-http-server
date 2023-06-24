package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sandromai/go-http-server/models"
	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

var LoginTokenTemplate string

func LoginTokenCreate(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "POST" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	var body *struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(request.Body).Decode(&body)

	if err == io.EOF {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "No data provided.",
		})

		return
	}

	if err != nil {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid data.",
		})

		return
	}

	body.Email = strings.TrimSpace(body.Email)

	if body.Email == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert your email address.",
		})

		return
	}

	if !utils.CheckEmail(body.Email) {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid email address.",
		})

		return
	}

	user, _ := (&models.User{}).FindByEmail(body.Email)

	if user != nil && user.Banned {
		utils.ReturnJSONResponse(writer, 403, &types.ReturnError{
			Error: "User banned.",
		})

		return
	}

	loginTokenModel := &models.LoginToken{}

	activeTokens, appErr := loginTokenModel.CountActiveByEmail(
		body.Email,
	)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	if activeTokens >= 3 {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "You've reached max active tokens, please wait to send new login tokens.",
		})

		return
	}

	lastTokenCreationTime, appErr := loginTokenModel.GetLastCreationTimeByEmail(
		body.Email,
	)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	if lastTokenCreationTime != "" {
		lastTokenTime, err := time.Parse(time.DateTime, lastTokenCreationTime)

		if err != nil {
			utils.ReturnJSONResponse(writer, 500, &types.ReturnError{
				Error: "Error parsing date.",
			})

			return
		}

		diff := time.Since(lastTokenTime)

		if diff.Seconds() < 60 {
			utils.ReturnJSONResponse(writer, 500, &types.ReturnError{
				Error: "Wait at least one minute to try again.",
			})

			return
		}
	}

	ipAddress := strings.Split(request.RemoteAddr, ":")[0]
	platform, browser := utils.GetDeviceInfo(request.Header.Get("User-Agent"))
	device := platform + ":" + browser
	expiresIn := int64(10 * 60)

	loginTokenId, appErr := loginTokenModel.Create(
		body.Email,
		ipAddress,
		device,
		expiresIn,
	)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	expiredAt := time.Now().Add(time.Minute * 10).Unix()

	loginToken, appErr := (&utils.JWT{}).Create(&types.LoginTokenPayload{
		LoginTokenId: loginTokenId,
		ExpiresAt:    expiredAt,
		CreatedAt:    time.Now().Unix(),
	})

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

	confirmAuthLink := "https://" + request.Host + "/auth/confirm?loginToken=" + loginToken

	emailBody := utils.UseTemplate(LoginTokenTemplate, map[string]string{"ConfirmAuthLink": confirmAuthLink})

	appErr = (&utils.Mailer{
		Host:     emailSettings.Host,
		Port:     emailSettings.Port,
		Username: emailSettings.Username,
		Password: emailSettings.Password,
	}).Send(
		"contact@company.com",
		"Company",
		body.Email,
		"",
		"Log in to Company Website",
		emailBody,
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
			LoginTokenId string `json:"loginTokenId"`
		}{LoginTokenId: loginTokenId},
	)
}
