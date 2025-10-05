package user

import (
	"context"
	"errors"

	"github.com/meehighlov/grats/internal/repositories/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) Save(ctx context.Context, user *entities.User) error {
	db, ok := ctx.Value(r.cfg.TxKey).(*gorm.DB)
	if !ok {
		return errors.New("not found transaction in context")
	}

	db = db.Session(&gorm.Session{
		SkipHooks: true,
	})

	_, _, _ = user.RefresTimestamps(r.cfg.Timezone)

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tg_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "tg_username", "chat_id", "is_admin", "updated_at"}),
	}).Create(user)
	if result.Error != nil {
		r.logger.Error("Error when trying to save user: " + result.Error.Error())
		return result.Error
	}

	r.logger.Debug("User created/updated")
	return nil
}
