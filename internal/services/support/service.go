package support

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/constants"
	"github.com/meehighlov/grats/internal/pagination"
	"github.com/meehighlov/grats/internal/repositories"
)

type Service struct {
	cfg          *config.Config
	logger       *slog.Logger
	repositories *repositories.Repositories
	clients      *clients.Clients
	builders     *builders.Builders
	constants    *constants.Constants
	pagination   *pagination.Pagination
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	repositories *repositories.Repositories,
	clients *clients.Clients,
	builders *builders.Builders,
	constants *constants.Constants,
	pagination *pagination.Pagination,
) *Service {
	return &Service{
		cfg:          cfg,
		logger:       logger,
		repositories: repositories,
		clients:      clients,
		builders:     builders,
		constants:    constants,
		pagination:   pagination,
	}
}
