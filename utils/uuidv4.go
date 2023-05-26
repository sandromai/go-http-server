package utils

import (
	"crypto/rand"
	"fmt"

	"github.com/sandromai/go-http-server/types"
)

func GenerateUUIDv4() (string, *types.AppError) {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to create random bytes.",
		}
	}

	bytes[6] = bytes[6]&15 | 64
	bytes[8] = bytes[8]&63 | 128

	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		bytes[:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:],
	), nil
}
