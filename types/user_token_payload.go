package types

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"time"
)

type UserTokenPayload struct {
	UserTokenId string `json:"userTokenId"`
	ExpiresAt   int64  `json:"expiresAt"`
	CreatedAt   int64  `json:"createdAt"`
}

func (payload *UserTokenPayload) ToJWT() (
	token string,
	appErr *AppError,
) {
	jsonHeaders, err := json.Marshal(&map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})

	if err != nil {
		return "", &AppError{
			StatusCode: 500,
			Message:    "Failed to create token headers.",
		}
	}

	encodedHeaders := base64.RawURLEncoding.EncodeToString([]byte(jsonHeaders))

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return "", &AppError{
			StatusCode: 500,
			Message:    "Failed to create token payload.",
		}
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(jsonPayload)

	hash := hmac.New(sha256.New, []byte(os.Getenv("JWT_KEY")))

	if _, err = hash.Write([]byte(encodedHeaders + "." + encodedPayload)); err != nil {
		return "", &AppError{
			StatusCode: 500,
			Message:    "Failed to create token signature.",
		}
	}

	encodedSignature := hex.EncodeToString(hash.Sum(nil))

	return encodedHeaders + "." + encodedPayload + "." + encodedSignature, nil
}

func (payload *UserTokenPayload) FromJWT(
	token string,
) *AppError {
	tokenParts := strings.Split(token, ".")

	payloadData, err := base64.RawURLEncoding.DecodeString(tokenParts[1])

	if err != nil {
		return &AppError{
			StatusCode: 500,
			Message:    "Failed to decode token payload.",
		}
	}

	err = json.Unmarshal(payloadData, payload)

	if err != nil {
		return &AppError{
			StatusCode: 500,
			Message:    "Failed to decode token payload data.",
		}
	}

	if payload.CreatedAt > time.Now().Unix() {
		return &AppError{
			StatusCode: 401,
			Message:    "Invalid token date.",
		}
	}

	if payload.ExpiresAt <= time.Now().Unix() {
		return &AppError{
			StatusCode: 401,
			Message:    "Expired token.",
		}
	}

	hash := hmac.New(sha256.New, []byte(os.Getenv("JWT_KEY")))

	if _, err = hash.Write([]byte(tokenParts[0] + "." + tokenParts[1])); err != nil {
		return &AppError{
			StatusCode: 500,
			Message:    "Failed to create token signature.",
		}
	}

	encodedSignature := hex.EncodeToString(hash.Sum(nil))

	if encodedSignature != tokenParts[2] {
		return &AppError{
			StatusCode: 401,
			Message:    "Invalid token.",
		}
	}

	return nil
}
