package types

type Admin struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	CreatedBy *int64 `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
}
