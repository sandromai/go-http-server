package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sandromai/go-http-server/middlewares"
	"github.com/sandromai/go-http-server/models"
	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

type Admin struct{}

func (*Admin) Login(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "POST" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	var body *struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		RememberMe bool   `json:"rememberMe"`
	}

	err := json.NewDecoder(request.Body).Decode(&body)

	if err == io.EOF {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert your username.",
		})

		return
	}

	if err != nil {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid data.",
		})

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

	var expiredAt int64

	if body.RememberMe {
		expiredAt = time.Now().Add(time.Hour * 24 * 30).Unix()
	} else {
		expiredAt = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	adminToken, appErr := (&types.AdminTokenPayload{
		AdminId:   admin.Id,
		ExpiresAt: expiredAt,
		CreatedAt: time.Now().Unix(),
	}).ToJWT()

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

func (*Admin) Register(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "POST" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	admin, appErr := middlewares.AuthenticateAdmin(request)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	var body *struct {
		Name            string `json:"name"`
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	err := json.NewDecoder(request.Body).Decode(&body)

	if err == io.EOF {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert the name.",
		})

		return
	}

	if err != nil {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid data.",
		})

		return
	}

	if body.Name == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert the name.",
		})

		return
	}

	if body.Username == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert the username.",
		})

		return
	}

	if body.Password == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert the password.",
		})

		return
	}

	if body.ConfirmPassword == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Repeat the password.",
		})

		return
	}

	if body.Password != body.ConfirmPassword {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "The passwords don't match.",
		})

		return
	}

	adminModel := &models.Admin{}

	adminId, appErr := adminModel.Create(
		body.Name,
		body.Username,
		body.Password,
		&admin.Id,
	)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	createdAdmin, appErr := adminModel.FindById(adminId)

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
		201,
		createdAdmin,
	)
}

func (*Admin) Update(
	writer http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != "PUT" {
		utils.ReturnJSONResponse(writer, 405, nil)

		return
	}

	admin, appErr := middlewares.AuthenticateAdmin(request)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	var body *struct {
		Name            string `json:"name"`
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	err := json.NewDecoder(request.Body).Decode(&body)

	if err == io.EOF {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert the name.",
		})

		return
	}

	if err != nil {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Invalid data.",
		})

		return
	}

	if body.Name == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert the name.",
		})

		return
	}

	if body.Username == "" {
		utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
			Error: "Insert the username.",
		})

		return
	}

	if body.Password != "" {
		if body.ConfirmPassword == "" {
			utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
				Error: "Repeat the password.",
			})

			return
		}

		if body.Password != body.ConfirmPassword {
			utils.ReturnJSONResponse(writer, 400, &types.ReturnError{
				Error: "The passwords don't match.",
			})

			return
		}
	}

	adminModel := &models.Admin{}

	appErr = adminModel.Update(
		body.Name,
		body.Username,
		body.Password,
		admin.Id,
	)

	if appErr != nil {
		utils.ReturnJSONResponse(
			writer,
			appErr.StatusCode,
			&types.ReturnError{Error: appErr.Message},
		)

		return
	}

	updatedAdmin, appErr := adminModel.FindById(admin.Id)

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
		updatedAdmin,
	)
}
