package wish_list

import (
	"context"

	"github.com/meehighlov/grats/internal/repositories/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ListFilter struct {
	WishListID string
	UserId     string
}

func (r *Repository) List(ctx context.Context, filter *ListFilter) ([]*models.WishList, error) {
	db, err := r.db.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	var wishLists []*models.WishList
	query := db.WithContext(ctx).Model(&models.WishList{})

	if filter.UserId != "" {
		query = query.Where("user_id = ?", filter.UserId)
	}
	if filter.WishListID != "" {
		query = query.Where("id = ?", filter.WishListID)
	}

	if err := query.Find(&wishLists).Error; err != nil {
		r.logger.Error("Error when filtering wish lists: " + err.Error())
		return nil, err
	}

	return wishLists, nil
}

func (r *Repository) Save(ctx context.Context, w *models.WishList) error {
	db, err := r.db.GetTx(ctx)
	if err != nil {
		return err
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

func (r *Repository) Delete(ctx context.Context, w *models.WishList) error {
	db, err := r.db.GetTx(ctx)
	if err != nil {
		return err
	}

	query := db.Model(&models.WishList{})

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

	result := query.Delete(&models.WishList{})
	if result.Error != nil {
		r.logger.Error("Error when trying to delete wishList: " + result.Error.Error())
		return result.Error
	}

	r.logger.Debug("WishList deleted")
	return nil
}
