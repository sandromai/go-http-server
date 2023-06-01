package types

type AdminTokenPayload struct {
	AdminId   string `json:"adminId"`
	ExpiresAt int64  `json:"expiresAt"`
	CreatedAt int64  `json:"createdAt"`
}
