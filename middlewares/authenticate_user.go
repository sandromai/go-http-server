package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/sandromai/go-http-server/models"
	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

func AuthenticateUser(
	request *http.Request,
	timezone *time.Location,
) (
	user *types.User,
	token string,
	appErr *types.AppError,
) {
	var userTokenId string

	userTokenModel := &models.UserToken{}
	userModel := &models.User{}

	ipAddress := strings.Split(request.RemoteAddr, ":")[0]
	platform, browser := utils.GetDeviceInfo(request.Header.Get("User-Agent"))

	var device string

	if platform != "" && browser != "" {
		device = platform + ":" + browser
	}

	loginTokenIdHeader := request.Header.Get("X-Login-Token-Id")

	if loginTokenIdHeader != "" {
		loginToken, appErr := (&models.LoginToken{}).FindById(
			loginTokenIdHeader,
		)

		if appErr != nil {
			return nil, "", appErr
		}

		tokenExpiresAt, err := time.ParseInLocation(
			time.DateTime,
			loginToken.ExpiresAt,
			timezone,
		)

		if err != nil {
			return nil, "", &types.AppError{
				StatusCode: 500,
				Message:    "Error parsing date.",
			}
		}

		if tokenExpiresAt.Before(time.Now()) {
			return nil, "", &types.AppError{
				StatusCode: 400,
				Message:    "Login token has expired.",
			}
		}

		tokenCreatedAt, err := time.ParseInLocation(
			time.DateTime,
			loginToken.CreatedAt,
			timezone,
		)

		if err != nil {
			return nil, "", &types.AppError{
				StatusCode: 500,
				Message:    "Error parsing date.",
			}
		}

		if tokenCreatedAt.After(time.Now()) {
			return nil, "", &types.AppError{
				StatusCode: 400,
				Message:    "Invalid login token date.",
			}
		}

		if loginToken.Denied {
			return nil, "", &types.AppError{
				StatusCode: 400,
				Message:    "Login token denied.",
			}
		}

		if !loginToken.Authorized {
			return nil, "", &types.AppError{
				StatusCode: 400,
				Message:    "Login token not authorized.",
			}
		}

		emailAvailable, appErr := userModel.CheckEmailAvailability(
			loginToken.Email,
		)

		if appErr != nil {
			return nil, "", appErr
		}

		if emailAvailable {
			userId, appErr := userModel.Create(
				loginToken.Email,
			)

			if appErr != nil {
				return nil, "", appErr
			}

			user, appErr = userModel.FindById(
				userId,
			)

			if appErr != nil {
				return nil, "", appErr
			}
		} else {
			user, appErr = userModel.FindByEmail(
				loginToken.Email,
			)

			if appErr != nil {
				return nil, "", appErr
			}
		}

		expiresIn := int64(30 * 24 * 60 * 60)

		userTokenId, appErr = userTokenModel.Create(
			user.Id,
			&loginToken.Id,
			nil,
			ipAddress,
			device,
			expiresIn,
		)

		if appErr != nil {
			return nil, "", appErr
		}

		expiredAt := time.Now().Add(30 * 24 * time.Hour).Unix()

		token, appErr = (&types.UserTokenPayload{
			UserTokenId: userTokenId,
			ExpiresAt:   expiredAt,
			CreatedAt:   time.Now().Unix(),
		}).ToJWT()

		if appErr != nil {
			return nil, "", appErr
		}
	} else {
		authorizationHeader := request.Header.Get("Authorization")

		if authorizationHeader == "" {
			return nil, "", &types.AppError{
				StatusCode: 401,
				Message:    "No authorization provided.",
			}
		}

		tokenParts := strings.Split(authorizationHeader, " ")

		if len(tokenParts) < 2 || tokenParts[0] != "Bearer" {
			return nil, "", &types.AppError{
				StatusCode: 401,
				Message:    "Invalid token.",
			}
		}

		userTokenPayload := &types.UserTokenPayload{}

		appErr := userTokenPayload.FromJWT(tokenParts[1])

		if appErr != nil {
			return nil, "", appErr
		}

		userTokenId = userTokenPayload.UserTokenId
	}

	userToken, appErr := userTokenModel.FindById(
		userTokenId,
	)

	if appErr != nil {
		return nil, "", appErr
	}

	tokenExpiresAt, err := time.ParseInLocation(
		time.DateTime,
		userToken.ExpiresAt,
		timezone,
	)

	if err != nil {
		return nil, "", &types.AppError{
			StatusCode: 500,
			Message:    "Error parsing date.",
		}
	}

	if tokenExpiresAt.Before(time.Now()) {
		tokenLastActivity, err := time.ParseInLocation(
			time.DateTime,
			userToken.LastActivity,
			timezone,
		)

		if err != nil {
			return nil, "", &types.AppError{
				StatusCode: 500,
				Message:    "Error parsing date.",
			}
		}

		timeGap := -(3 * 24) * time.Hour

		if tokenLastActivity.After(time.Now().Add(timeGap)) && tokenExpiresAt.After(time.Now().Add(timeGap)) {
			expiresIn := int64(30 * 24 * 60 * 60)

			userTokenId, appErr = userTokenModel.Create(
				userToken.Id,
				nil,
				&userToken.Id,
				ipAddress,
				device,
				expiresIn,
			)

			if appErr != nil {
				return nil, "", appErr
			}

			expiredAt := time.Now().Add(30 * 24 * time.Hour).Unix()

			token, appErr = (&types.UserTokenPayload{
				UserTokenId: userTokenId,
				ExpiresAt:   expiredAt,
				CreatedAt:   time.Now().Unix(),
			}).ToJWT()

			if appErr != nil {
				return nil, "", appErr
			}

			userToken, appErr = userTokenModel.FindById(
				userTokenId,
			)

			if appErr != nil {
				return nil, "", appErr
			}
		} else {
			return nil, "", &types.AppError{
				StatusCode: 400,
				Message:    "User token has expired.",
			}
		}
	}

	tokenCreatedAt, err := time.ParseInLocation(
		time.DateTime,
		userToken.CreatedAt,
		timezone,
	)

	if err != nil {
		return nil, "", &types.AppError{
			StatusCode: 500,
			Message:    "Error parsing date.",
		}
	}

	if tokenCreatedAt.After(time.Now()) {
		return nil, "", &types.AppError{
			StatusCode: 400,
			Message:    "Invalid user token date.",
		}
	}

	if userToken.Disconnected {
		return nil, "", &types.AppError{
			StatusCode: 400,
			Message:    "Session disconnected.",
		}
	}

	if user == nil {
		user, appErr = userModel.FindById(
			userToken.UserId,
		)

		if appErr != nil {
			return nil, "", appErr
		}
	}

	if user.Banned {
		return nil, "", &types.AppError{
			StatusCode: 403,
			Message:    "User banned.",
		}
	}

	userTokenModel.UpdateActivity(
		userToken.Id,
	)

	return user, token, nil
}
