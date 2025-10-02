package services

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/builders"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/constants"
	"github.com/meehighlov/grats/internal/pagination"
	"github.com/meehighlov/grats/internal/repositories"
	"github.com/meehighlov/grats/internal/services/support"
	"github.com/meehighlov/grats/internal/services/user"
	"github.com/meehighlov/grats/internal/services/wish"
	"github.com/meehighlov/grats/internal/services/wishlist"
)

type Services struct {
	User     *user.Service
	Wish     *wish.Service
	WishList *wishlist.Service
	Support  *support.Service
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	repositories *repositories.Repositories,
	clients *clients.Clients,
	builders *builders.Builders,
	constants *constants.Constants,
	pagination *pagination.Pagination,
) *Services {
	return &Services{
		User:     user.New(cfg, logger, repositories, clients, builders, constants, pagination),
		Wish:     wish.New(cfg, logger, repositories, clients, builders, constants, pagination),
		WishList: wishlist.New(cfg, logger, repositories, clients, builders, constants, pagination),
		Support:  support.New(cfg, logger, repositories, clients, builders, constants, pagination),
	}
}
