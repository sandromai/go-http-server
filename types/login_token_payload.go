package types

type LoginTokenPayload struct {
	LoginTokenId string `json:"loginTokenId"`
	ExpiresAt    int64  `json:"expiresAt"`
	CreatedAt    int64  `json:"createdAt"`
}
