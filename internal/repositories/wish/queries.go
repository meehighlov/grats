package wish

import (
	"context"

	"github.com/meehighlov/grats/internal/repositories/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ListFilter struct {
	Limit      int
	Offset     int
	ChatID     string
	WishListID string
	UserId     string
}

type CountFilter struct {
	WishListID string
}

func (r *Repository) List(ctx context.Context, filter *ListFilter) ([]*models.Wish, error) {
	db, err := r.db.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	var wishes []*models.Wish
	query := db.Model(&models.Wish{})
	if filter.Limit != 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset != 0 {
		query = query.Offset(filter.Offset)
	}
	if filter.WishListID != "" {
		query = query.Where("wish_list_id = ?", filter.WishListID)
	}
	if filter.UserId != "" {
		query = query.Where("user_id = ?", filter.UserId)
	}
	query = query.Order("executor_id DESC")

	if err := query.Find(&wishes).Error; err != nil {
		r.logger.Error("Error when list wishes: " + err.Error())
		return nil, err
	}

	return wishes, nil
}

func (r *Repository) Count(ctx context.Context, filter *CountFilter) (int64, error) {
	db, err := r.db.GetTx(ctx)
	if err != nil {
		return 0, err
	}

	var count int64
	query := db.Model(&models.Wish{})

	if filter.WishListID != "" {
		query = query.Where("wish_list_id = ?", filter.WishListID)
	}

	if err := query.Count(&count).Error; err != nil {
		r.logger.Error("Error when counting wishes: " + err.Error())
		return 0, err
	}

	return count, nil
}

func (r *Repository) GetWithLock(ctx context.Context, w *models.Wish) ([]*models.Wish, error) {
	db, err := r.db.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	var wishes []*models.Wish
	query := db.Model(&models.Wish{}).Clauses(clause.Locking{Strength: "UPDATE"})

	if w.ID != "" {
		query = query.Where("id = ?", w.ID)
	}

	if err := query.Find(&wishes).Error; err != nil {
		r.logger.Error("Error when getting wishes with lock: " + err.Error())
		return nil, err
	}

	return wishes, nil
}

func (r *Repository) Get(ctx context.Context, wishId string) (*models.Wish, error) {
	db, err := r.db.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	var wish models.Wish
	query := db.Model(&models.Wish{}).Where("id = ?", wishId)

	if err := query.First(&wish).Error; err != nil {
		r.logger.Error("Error when getting wish: " + err.Error())
		return nil, err
	}

	return &wish, nil
}

func (r *Repository) Save(ctx context.Context, w *models.Wish) error {
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
		DoUpdates: clause.AssignmentColumns([]string{"name", "chat_id", "user_id", "wish_list_id", "link", "price", "executor_id", "updated_at"}),
	}).Create(w)
	if result.Error != nil {
		r.logger.Error("Error when trying to save wish: " + result.Error.Error())
		return result.Error
	}

	r.logger.Debug("Wish created/updated")
	return nil
}

func (r *Repository) Delete(ctx context.Context, w *models.Wish) error {
	db, err := r.db.GetTx(ctx)
	if err != nil {
		return err
	}

	db = db.Session(&gorm.Session{
		SkipHooks: true,
	})

	query := db.Model(&models.Wish{})

	if w.ID != "" {
		query = query.Where("id = ?", w.ID)
	}
	if w.UserId != "" {
		query = query.Where("user_id = ?", w.UserId)
	}
	if w.ChatId != "" {
		query = query.Where("chat_id = ?", w.ChatId)
	}
	if w.WishListId != "" {
		query = query.Where("wish_list_id = ?", w.WishListId)
	}

	result := query.Delete(&models.Wish{})
	if result.Error != nil {
		r.logger.Error("Error when trying to delete wish: " + result.Error.Error())
		return result.Error
	}

	r.logger.Debug("Wish deleted")
	return nil
}
