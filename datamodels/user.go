package datamodels

type User struct {
	ID        int64  `json:"id" form:"id"`
	Nickname  string `json:"nickname" form:"nickname"`
	Username  string `json:"username" form:"username"`
	CreatedAt int64  `json:"created_at" form:"-"`
	UpdatedAt int64  `json:"update_at" form:"-"`
}
