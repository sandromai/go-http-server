package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/sandromai/go-http-server/models"
	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

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

	if strings.TrimSpace(body.Email) == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert your email address.",
		})

		return
	}

	if matched, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, strings.TrimSpace(body.Email)); !matched || err != nil {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid email address.",
		})

		return
	}

	user, _ := (&models.User{}).FindByEmail(body.Email)

	if user != nil && user.Banned {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "User banned.",
		})

		return
	}

	// check quantity of active login tokens
	// check if last login token was created more than 1 minute ago

	ipAddress := strings.Split(request.RemoteAddr, ":")[0]
	platform, browser := utils.GetDeviceInfo(request.Header.Get("User-Agent"))
	device := platform + ":" + browser
	expiresIn := int64(10 * 60)

	loginTokenId, appErr := (&models.LoginToken{}).Create(
		strings.TrimSpace(body.Email),
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

	// create link user loginToken
	// send link to email
	fmt.Println(loginToken)

	utils.ReturnJSONResponse(
		writer,
		200,
		&struct {
			LoginTokenId string `json:"loginTokenId"`
		}{LoginTokenId: loginTokenId},
	)
}
