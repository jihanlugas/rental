package request

type Signin struct {
	Username string `db:"username,use_zero" json:"username" form:"username" query:"username" validate:"required"`
	Passwd   string `db:"passwd,use_zero" json:"passwd" form:"passwd" query:"passwd" validate:"required,lte=200"`
}

type ChangePassword struct {
	CurrentPasswd string `json:"currentPasswd" form:"currentPasswd" validate:"required,lte=200"`
	Passwd        string `json:"passwd" form:"passwd" validate:"required,lte=200"`
	ConfirmPasswd string `json:"confirmPasswd" form:"confirmPasswd" validate:"required,lte=200,eqfield=Passwd"`
}
