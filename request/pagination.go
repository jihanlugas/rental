package request

import (
	"github.com/jihanlugas/rental/config"
)

type Paging struct {
	Page  int `json:"page,omitempty" form:"page" query:"page"`
	Limit int `json:"limit,omitempty" form:"limit" query:"limit"`
}

func (p *Paging) GetPage() int {
	return p.Page
}

func (p *Paging) GetLimit() int {
	if p.Limit <= config.DataPerPage {
		return config.DataPerPage
	} else {
		return p.Limit
	}
}

func (p *Paging) SetLimit(lim int) {
	p.Limit = lim
}

func (p *Paging) SetPage(page int) {
	p.Page = page
}

type IPaging interface {
	GetPage() int
	GetLimit() int
	SetLimit(int)
	SetPage(int)
}
