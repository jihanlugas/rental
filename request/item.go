package request

type CreateItem struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateItem struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type PageItem struct {
	Paging
	Name        string `db:"name,use_zero" json:"name" form:"name" query:"name" validate:"lte=200"`
	Description string `db:"description,use_zero" json:"description" form:"description" query:"description" validate:"lte=200"`
}

type ListItem struct {
	CompanyID   string `db:"company_id,use_zero" json:"companyId" form:"companyId" query:"companyId" validate:"lte=200"`
	Name        string `db:"name,use_zero" json:"name" form:"name" query:"name" validate:"lte=200"`
	Description string `db:"description,use_zero" json:"description" form:"description" query:"description" validate:"lte=200"`
}
