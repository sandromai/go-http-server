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

type LoginToken struct {
	Template string
	Timezone *time.Location
}

func (l *LoginToken) Create(
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
			Error: "Insert your email address.",
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
		lastTokenTime, err := time.ParseInLocation(
			time.DateTime,
			lastTokenCreationTime,
			l.Timezone,
		)

		if err != nil {
			utils.ReturnJSONResponse(writer, 500, &types.ReturnError{
				Error: "Error parsing date.",
			})

			return
		}

		if time.Since(lastTokenTime).Seconds() < 60 {
			utils.ReturnJSONResponse(writer, 500, &types.ReturnError{
				Error: "Wait at least one minute to try again.",
			})

			return
		}
	}

	ipAddress := strings.Split(request.RemoteAddr, ":")[0]
	platform, browser := utils.GetDeviceInfo(request.Header.Get("User-Agent"))

	var device string

	if platform != "" && browser != "" {
		device = platform + ":" + browser
	}

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

	loginTokenString, appErr := (&types.LoginTokenPayload{
		LoginTokenId: loginTokenId,
		ExpiresAt:    expiredAt,
		CreatedAt:    time.Now().Unix(),
	}).ToJWT()

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

	confirmAuthLink := "https://" + request.Host + "/auth/confirm?loginToken=" + loginTokenString

	emailBody := utils.UseTemplate(l.Template, map[string]string{"ConfirmAuthLink": confirmAuthLink})

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

func (l *LoginToken) Check(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "POST" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	var body *struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(request.Body).Decode(&body)

	if err == io.EOF {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Login token not identified.",
		})

		return
	}

	if err != nil {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid data.",
		})

		return
	}

	body.Token = strings.TrimSpace(body.Token)

	if body.Token == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Login token not identified.",
		})

		return
	}

	loginTokenPayload := &types.LoginTokenPayload{}

	appErr := loginTokenPayload.FromJWT(body.Token)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	loginToken, appErr := (&models.LoginToken{}).FindById(
		loginTokenPayload.LoginTokenId,
	)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	tokenExpiresAt, err := time.ParseInLocation(
		time.DateTime,
		loginToken.ExpiresAt,
		l.Timezone,
	)

	if err != nil {
		utils.ReturnJSONResponse(writer, 500, &types.ReturnError{
			Error: "Error parsing date.",
		})

		return
	}

	if tokenExpiresAt.Before(time.Now()) {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Login token has expired.",
		})

		return
	}

	tokenCreatedAt, err := time.ParseInLocation(
		time.DateTime,
		loginToken.CreatedAt,
		l.Timezone,
	)

	if err != nil {
		utils.ReturnJSONResponse(writer, 500, &types.ReturnError{
			Error: "Error parsing date.",
		})

		return
	}

	if tokenCreatedAt.After(time.Now()) {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid login token date.",
		})

		return
	}

	if loginToken.Denied {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "This login token was denied.",
		})

		return
	}

	if loginToken.Authorized {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "This login token was already authorized.",
		})

		return
	}

	utils.ReturnJSONResponse(
		writer,
		200,
		nil,
	)
}

func (l *LoginToken) Deny(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "POST" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	var body *struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(request.Body).Decode(&body)

	if err == io.EOF {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Login token not identified.",
		})

		return
	}

	if err != nil {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid data.",
		})

		return
	}

	body.Token = strings.TrimSpace(body.Token)

	if body.Token == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Login token not identified.",
		})

		return
	}

	loginTokenPayload := &types.LoginTokenPayload{}

	appErr := loginTokenPayload.FromJWT(body.Token)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	loginTokenModel := &models.LoginToken{}

	loginToken, appErr := loginTokenModel.FindById(
		loginTokenPayload.LoginTokenId,
	)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	tokenExpiresAt, err := time.ParseInLocation(
		time.DateTime,
		loginToken.ExpiresAt,
		l.Timezone,
	)

	if err != nil {
		utils.ReturnJSONResponse(writer, 500, &types.ReturnError{
			Error: "Error parsing date.",
		})

		return
	}

	if tokenExpiresAt.Before(time.Now()) {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Login token has expired.",
		})

		return
	}

	tokenCreatedAt, err := time.ParseInLocation(
		time.DateTime,
		loginToken.CreatedAt,
		l.Timezone,
	)

	if err != nil {
		utils.ReturnJSONResponse(writer, 500, &types.ReturnError{
			Error: "Error parsing date.",
		})

		return
	}

	if tokenCreatedAt.After(time.Now()) {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid login token date.",
		})

		return
	}

	if loginToken.Authorized {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "This login token was already authorized.",
		})

		return
	}

	if loginToken.Denied {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "This login token was already denied.",
		})

		return
	}

	appErr = loginTokenModel.Deny(loginToken.Id)

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

func (l *LoginToken) Authorize(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "POST" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	var body *struct {
		Token string `json:"token"`
	}

	err := json.NewDecoder(request.Body).Decode(&body)

	if err == io.EOF {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Login token not identified.",
		})

		return
	}

	if err != nil {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid data.",
		})

		return
	}

	body.Token = strings.TrimSpace(body.Token)

	if body.Token == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Login token not identified.",
		})

		return
	}

	loginTokenPayload := &types.LoginTokenPayload{}

	appErr := loginTokenPayload.FromJWT(body.Token)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	loginTokenModel := &models.LoginToken{}

	loginToken, appErr := loginTokenModel.FindById(
		loginTokenPayload.LoginTokenId,
	)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	tokenExpiresAt, err := time.ParseInLocation(
		time.DateTime,
		loginToken.ExpiresAt,
		l.Timezone,
	)

	if err != nil {
		utils.ReturnJSONResponse(writer, 500, &types.ReturnError{
			Error: "Error parsing date.",
		})

		return
	}

	if tokenExpiresAt.Before(time.Now()) {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Login token has expired.",
		})

		return
	}

	tokenCreatedAt, err := time.ParseInLocation(
		time.DateTime,
		loginToken.CreatedAt,
		l.Timezone,
	)

	if err != nil {
		utils.ReturnJSONResponse(writer, 500, &types.ReturnError{
			Error: "Error parsing date.",
		})

		return
	}

	if tokenCreatedAt.After(time.Now()) {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid login token date.",
		})

		return
	}

	if loginToken.Authorized {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "This login token was already authorized.",
		})

		return
	}

	if loginToken.Denied {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "This login token was already denied.",
		})

		return
	}

	appErr = loginTokenModel.Authorize(loginToken.Id)

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
