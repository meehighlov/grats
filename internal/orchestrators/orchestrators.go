package orchestrators

import (
	"log/slog"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/orchestrators/support"
	"github.com/meehighlov/grats/internal/orchestrators/user"
	"github.com/meehighlov/grats/internal/orchestrators/wish"
	"github.com/meehighlov/grats/internal/orchestrators/wishlist"
	"github.com/meehighlov/grats/internal/services"
	"gorm.io/gorm"
)

type Orchestrators struct {
	User     *user.Orchestrator
	Wish     *wish.Orchestrator
	WishList *wishlist.Orchestrator
	Support  *support.Orchestrator
}

func New(cfg *config.Config, logger *slog.Logger, db *gorm.DB, services *services.Services) *Orchestrators {
	return &Orchestrators{
		User:     user.New(cfg, logger, services, db),
		Wish:     wish.New(cfg, logger, services, db),
		WishList: wishlist.New(cfg, logger, services, db),
		Support:  support.New(cfg, logger, services, db),
	}
}
