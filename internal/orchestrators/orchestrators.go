package orchestrators

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/orchestrators/support"
	"github.com/meehighlov/grats/internal/orchestrators/user"
	"github.com/meehighlov/grats/internal/orchestrators/wish"
	"github.com/meehighlov/grats/internal/orchestrators/wishlist"
	"github.com/meehighlov/grats/internal/services"
)

type Orchestrators struct {
	User     *user.Orchestrator
	Wish     *wish.Orchestrator
	WishList *wishlist.Orchestrator
	Support  *support.Orchestrator
}

func New(cfg *config.Config, logger *slog.Logger, services *services.Services) *Orchestrators {
	return &Orchestrators{
		User:     user.New(cfg, logger, services),
		Wish:     wish.New(cfg, logger, services),
		WishList: wishlist.New(cfg, logger, services),
		Support:  support.New(cfg, logger, services),
	}
}
