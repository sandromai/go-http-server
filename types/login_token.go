package types

type LoginToken struct {
	Id         int64  `json:"id"`
	Email      string `json:"email"`
	DeviceCode string `json:"deviceCode"`
	IPAddress  string `json:"IPAddress"`
	Device     string `json:"device"`
	Authorized bool   `json:"authorized"`
	Denied     bool   `json:"denied"`
	ExpiresAt  string `json:"expiresAt"`
	CreatedAt  string `json:"createdAt"`
}
