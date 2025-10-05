package repositories

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/repositories/user"
	"github.com/meehighlov/grats/internal/repositories/wish"
	"github.com/meehighlov/grats/internal/repositories/wish_list"
)

type Repositories struct {
	User     *user.Repository
	Wish     *wish.Repository
	WishList *wish_list.Repository
}

func New(cfg *config.Config, logger *slog.Logger) *Repositories {
	return &Repositories{
		User:     user.New(cfg, logger),
		Wish:     wish.New(cfg, logger),
		WishList: wish_list.New(cfg, logger),
	}
}
