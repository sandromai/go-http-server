package types

type Admin struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Username  string  `json:"username"`
	CreatedBy *string `json:"createdBy"`
	CreatedAt string  `json:"createdAt"`
}
