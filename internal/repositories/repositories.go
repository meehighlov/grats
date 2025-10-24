package repositories

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/infra/redis"
	"github.com/meehighlov/grats/internal/repositories/cache"
	"github.com/meehighlov/grats/internal/repositories/state"
	"github.com/meehighlov/grats/internal/repositories/user"
	"github.com/meehighlov/grats/internal/repositories/wish"
	"github.com/meehighlov/grats/internal/repositories/wish_list"
	"gorm.io/gorm"
)

type Repositories struct {
	User     *user.Repository
	Wish     *wish.Repository
	WishList *wish_list.Repository
	Cache    *cache.Repository
	State    *state.Repository
}

func New(cfg *config.Config, logger *slog.Logger, db *gorm.DB, redis *redis.Client) *Repositories {
	return &Repositories{
		User:     user.New(cfg, logger),
		Wish:     wish.New(cfg, logger),
		WishList: wish_list.New(cfg, logger),
		Cache:    cache.New(cfg, logger, redis),
		State:    state.New(cfg, logger, redis),
	}
}
