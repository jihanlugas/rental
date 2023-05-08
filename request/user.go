package request

type Signin struct {
	Username string `db:"username,use_zero" json:"username" form:"username" query:"username" validate:"required"`
	Passwd   string `db:"passwd,use_zero" json:"passwd" form:"passwd" query:"passwd" validate:"required,lte=200"`
}
