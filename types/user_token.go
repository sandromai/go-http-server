package types

type UserToken struct {
	Id             string  `json:"id"`
	UserId         string  `json:"userId"`
	FromLoginToken *string `json:"fromLoginToken"`
	FromUserToken  *string `json:"fromUserToken"`
	IPAddress      string  `json:"IPAddress"`
	Device         string  `json:"device"`
	Disconnected   bool    `json:"disconnected"`
	LastActivity   string  `json:"lastActivity"`
	ExpiresAt      string  `json:"expiresAt"`
	CreatedAt      string  `json:"createdAt"`
}
