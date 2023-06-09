package middlewares

import (
	"net/http"
	"strings"

	"github.com/sandromai/go-http-server/models"
	"github.com/sandromai/go-http-server/types"
)

func AuthenticateAdmin(
	request *http.Request,
) (
	*types.Admin,
	*types.AppError,
) {
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

	adminTokenPayload := &types.AdminTokenPayload{}

	appErr := adminTokenPayload.FromJWT(tokenParts[1])

	if appErr != nil {
		return nil, appErr
	}

	admin, appErr := (&models.Admin{}).FindById(
		adminTokenPayload.AdminId,
	)

	if appErr != nil {
		return nil, appErr
	}

	return admin, nil
}
