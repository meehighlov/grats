package pagination

import (
	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/config"
)

type Pagination struct {
	builders   *builders.Builders
	BaseOffset int
	Limit      int
}

func New(cfg *config.Config, builders *builders.Builders) *Pagination {
	return &Pagination{
		builders:   builders,
		BaseOffset: 5,
		Limit:      cfg.ListLimitLen,
	}
}
