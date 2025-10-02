package user

import (
	"context"

	"github.com/meehighlov/grats/internal/repositories/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) Save(ctx context.Context, tx *gorm.DB, user *entities.User) error {
	db := r.db
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
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
