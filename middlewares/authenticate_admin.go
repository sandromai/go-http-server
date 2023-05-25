package middlewares

import (
	"net/http"
	"strings"

	"github.com/sandromai/go-http-server/types"
	"github.com/sandromai/go-http-server/utils"
)

func AuthenticateAdmin(
	request *http.Request,
) (*types.Admin, *types.AppError) {
	authorizationHeader := request.Header.Get("Authorization")

	if authorizationHeader == "" {
		return nil, &types.AppError{
			StatusCode: 401,
			Message:    "No authorization provided.",
		}
	}

	tokenParts := strings.Split(authorizationHeader, " ")

	if len(tokenParts) < 2 || tokenParts[0] != "Bearer" {
		return nil, &types.AppError{
			StatusCode: 401,
			Message:    "Invalid token.",
		}
	}

	admin, appErr := utils.AuthenticateAdmin(tokenParts[1])

	if appErr != nil {
		return nil, appErr
	}

	return admin, nil
}
