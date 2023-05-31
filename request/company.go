package request

type CreateCompany struct {
	Name   string `json:"name" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}

type UpdateCompany struct {
	Name   string `json:"name" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}
