package wish_list

import (
	"context"
	"errors"

	"github.com/meehighlov/grats/internal/repositories/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) Filter(ctx context.Context, w *entities.WishList) ([]*entities.WishList, error) {
	db, ok := ctx.Value(r.cfg.TxKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("not found transaction in context")
	}

	var wishLists []*entities.WishList
	query := db.Model(&entities.WishList{})

	if w.UserId != "" {
		query = query.Where("user_id = ?", w.UserId)
	}
	if w.Name != "" {
		query = query.Where("name = ?", w.Name)
	}
	if w.ChatId != "" {
		query = query.Where("chat_id = ?", w.ChatId)
	}
	if w.ID != "" {
		query = query.Where("id = ?", w.ID)
	}

	if err := query.Find(&wishLists).Error; err != nil {
		r.logger.Error("Error when filtering wish lists: " + err.Error())
		return nil, err
	}

	return wishLists, nil
}

func (r *Repository) Save(ctx context.Context, w *entities.WishList) error {
	db, ok := ctx.Value(r.cfg.TxKey).(*gorm.DB)
	if !ok {
		return errors.New("not found transaction in context")
	}

	db = db.Session(&gorm.Session{
		SkipHooks: true,
	})

	_, _, _ = w.RefresTimestamps(r.cfg.Timezone)

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "user_id", "chat_id", "updated_at"}),
	}).Create(w)
	if result.Error != nil {
		r.logger.Error("Error when trying to save wishList: " + result.Error.Error())
		return result.Error
	}

	r.logger.Debug("WishList created/updated")
	return nil
}

func (r *Repository) Delete(ctx context.Context, w *entities.WishList) error {
	db, ok := ctx.Value(r.cfg.TxKey).(*gorm.DB)
	if !ok {
		return errors.New("not found transaction in context")
	}

	query := db.Model(&entities.WishList{})

	if w.ID != "" {
		query = query.Where("id = ?", w.ID)
	}
	if w.UserId != "" {
		query = query.Where("user_id = ?", w.UserId)
	}
	if w.ChatId != "" {
		query = query.Where("chat_id = ?", w.ChatId)
	}
	if w.Name != "" {
		query = query.Where("name = ?", w.Name)
	}

	result := query.Delete(&entities.WishList{})
	if result.Error != nil {
		r.logger.Error("Error when trying to delete wishList: " + result.Error.Error())
		return result.Error
	}

	r.logger.Debug("WishList deleted")
	return nil
}
