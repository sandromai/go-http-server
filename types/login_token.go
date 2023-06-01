package types

type LoginToken struct {
	Id         string `json:"id"`
	Email      string `json:"email"`
	IPAddress  string `json:"IPAddress"`
	Device     string `json:"device"`
	Authorized bool   `json:"authorized"`
	Denied     bool   `json:"denied"`
	ExpiresAt  string `json:"expiresAt"`
	CreatedAt  string `json:"createdAt"`
}
