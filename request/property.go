package request

type CreateProperty struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"descriptionId" validate:"required"`
}

type UpdateProperty struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"descriptionId" validate:"required"`
}
