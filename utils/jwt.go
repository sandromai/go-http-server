package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/sandromai/go-http-server/types"
)

type JWT struct{}

func (*JWT) Create(
	payload *types.AdminTokenPayload,
) (string, *types.AppError) {
	jsonHeaders, err := json.Marshal(&map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to create token headers.",
		}
	}

	encodedHeaders := base64.RawURLEncoding.EncodeToString([]byte(jsonHeaders))

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to create token payload.",
		}
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString([]byte(jsonPayload))

	hash := hmac.New(sha256.New, []byte(os.Getenv("JWT_KEY")))

	if _, err = hash.Write([]byte(encodedHeaders + "." + encodedPayload)); err != nil {
		return "", &types.AppError{
			StatusCode: 500,
			Message:    "Failed to create token signature.",
		}
	}

	encodedSignature := hex.EncodeToString(hash.Sum(nil))

	return encodedHeaders + "." + encodedPayload + "." + encodedSignature, nil
}

func (*JWT) Check(
	token string,
) (
	*types.AdminTokenPayload,
	*types.AppError,
) {
	tokenParts := strings.Split(token, ".")

	payload, err := base64.RawURLEncoding.DecodeString(tokenParts[1])

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to decode token payload.",
		}
	}

	var payloadData *types.AdminTokenPayload

	err = json.Unmarshal(payload, &payloadData)

	if err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to decode token payload data.",
		}
	}

	if payloadData.CreatedAt > time.Now().Unix() {
		return nil, &types.AppError{
			StatusCode: 401,
			Message:    "Invalid token date.",
		}
	}

	if payloadData.ExpiresAt <= time.Now().Unix() {
		return nil, &types.AppError{
			StatusCode: 401,
			Message:    "Expired token.",
		}
	}

	hash := hmac.New(sha256.New, []byte(os.Getenv("JWT_KEY")))

	if _, err = hash.Write([]byte(tokenParts[0] + "." + tokenParts[1])); err != nil {
		return nil, &types.AppError{
			StatusCode: 500,
			Message:    "Failed to create token signature.",
		}
	}

	encodedSignature := hex.EncodeToString(hash.Sum(nil))

	if encodedSignature != tokenParts[2] {
		return nil, &types.AppError{
			StatusCode: 401,
			Message:    "Invalid token.",
		}
	}

	return payloadData, nil
}
