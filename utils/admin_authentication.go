package utils

import (
	"github.com/sandromai/go-http-server/models"
	"github.com/sandromai/go-http-server/types"
)

func AuthenticateAdmin(
	adminToken string,
) (
	*types.Admin,
	*types.AppError,
) {
	adminTokenPayload, appErr := (&JWT{}).Check(adminToken)

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
