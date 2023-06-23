package types

type UserTokenPayload struct {
	UserTokenId string `json:"userTokenId"`
	ExpiresAt   int64  `json:"expiresAt"`
	CreatedAt   int64  `json:"createdAt"`
}
