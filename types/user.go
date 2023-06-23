package types

type User struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	Banned    bool   `json:"banned"`
	CreatedAt string `json:"createdAt"`
}
