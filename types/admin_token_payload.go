package types

type AdminTokenPayload struct {
	AdminId   int64 `json:"adminId"`
	ExpiresAt int64 `json:"expiresAt"`
	CreatedAt int64 `json:"createdAt"`
}
