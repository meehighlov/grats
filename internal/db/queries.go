package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/meehighlov/grats/internal/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (user *User) GetId() string {
	return user.ID
}

func (user *User) Save(ctx context.Context, tx *gorm.DB) error {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	db = db.Session(&gorm.Session{
		SkipHooks: true,
	})

	_, _, _ = user.RefresTimestamps()

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tg_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "tg_username", "chat_id", "is_admin", "updated_at"}),
	}).Create(user)
	if result.Error != nil {
		slog.Error("Error when trying to save user: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("User created/updated")
	return nil
}

func (u *User) GreaterThan(other common.PaginatedEntity) bool {
	return true
}

func (u *User) ButtonText() string {
	return u.Name
}

func (u *User) Search(ctx context.Context, tx *gorm.DB, params *common.SearchParams) ([]common.PaginatedEntity, error) {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	var users []*User
	query := db.Model(&User{})

	if params.ListId != "" {
		query = query.Where("id = ?", params.ListId)
	}

	if err := query.Find(&users).Error; err != nil {
		slog.Error("Error when searching users: " + err.Error())
		return nil, err
	}

	var entities []common.PaginatedEntity
	for _, user := range users {
		entities = append(entities, user)
	}

	return entities, nil
}

func (wish *Wish) GetId() string {
	return wish.ID
}

func (w *Wish) GreaterThan(other common.PaginatedEntity) bool {
	otherWish, ok := other.(*Wish)
	if !ok {
		return false
	}

	if w.Price == "" && otherWish.Price != "" {
		return false
	}
	if w.Price != "" && otherWish.Price == "" {
		return true
	}
	if w.Price == "" && otherWish.Price == "" {
		return w.Name > otherWish.Name
	}

	// parsing price
	// w.Price and otherWish.Price

	return w.Price > otherWish.Price
}

func (w *Wish) ButtonText() string {
	price := ""
	if w.Price != "" {
		price = fmt.Sprintf(" - %s(RUB)", w.Price)
	}
	buttonText := fmt.Sprintf("%s%s", w.Name, price)

	return buttonText
}

func (w *Wish) Search(ctx context.Context, tx *gorm.DB, params *common.SearchParams) ([]common.PaginatedEntity, error) {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	var wishes []*Wish
	query := db.Model(&Wish{})

	if params.ListId != "" {
		query = query.Where("wish_list_id = ?", params.ListId)
	}

	if err := query.Find(&wishes).Error; err != nil {
		slog.Error("Error when searching wishes: " + err.Error())
		return nil, err
	}

	var entities []common.PaginatedEntity
	for _, wish := range wishes {
		entities = append(entities, wish)
	}

	return entities, nil
}

func (w *Wish) Filter(ctx context.Context, tx *gorm.DB) ([]*Wish, error) {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	var wishes []*Wish
	query := db.Model(&Wish{})

	if w.UserId != "" {
		query = query.Where("user_id = ?", w.UserId)
	}
	if w.ID != "" {
		query = query.Where("id = ?", w.ID)
	}
	if w.ChatId != "" {
		query = query.Where("chat_id = ?", w.ChatId)
	}
	if w.WishListId != "" {
		query = query.Where("wish_list_id = ?", w.WishListId)
	}

	if err := query.Find(&wishes).Error; err != nil {
		slog.Error("Error when filtering wishes: " + err.Error())
		return nil, err
	}

	return wishes, nil
}

func (w *Wish) GetWithLock(ctx context.Context, tx *gorm.DB) ([]*Wish, error) {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	var wishes []*Wish
	query := db.Model(&Wish{})

	if w.UserId != "" {
		query = query.Where("user_id = ?", w.UserId)
	}
	if w.ID != "" {
		query = query.Where("id = ?", w.ID)
	}
	if w.ChatId != "" {
		query = query.Where("chat_id = ?", w.ChatId)
	}
	if w.WishListId != "" {
		query = query.Where("wish_list_id = ?", w.WishListId)
	}

	query = query.Set("gorm:query_option", "FOR UPDATE")

	if err := query.Find(&wishes).Error; err != nil {
		slog.Error("Error when getting wishes with lock: " + err.Error())
		return nil, err
	}

	return wishes, nil
}

func (w *Wish) Save(ctx context.Context, tx *gorm.DB) error {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	db = db.Session(&gorm.Session{
		SkipHooks: true,
	})

	_, _, _ = w.RefresTimestamps()

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "chat_id", "user_id", "link", "executor_id", "price", "wish_list_id", "updated_at"}),
	}).Create(w)
	if result.Error != nil {
		slog.Error("Error when trying to save wish: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("Wish created/updated")
	return nil
}

func (w *Wish) Delete(ctx context.Context, tx *gorm.DB) error {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	query := db.Model(&Wish{})

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

	result := query.Delete(&Wish{})
	if result.Error != nil {
		slog.Error("Error when trying to delete wish: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("Wish deleted")
	return nil
}

func (wishList *WishList) GetId() string {
	return wishList.ID
}

func (w *WishList) GreaterThan(other common.PaginatedEntity) bool {
	return true
}

func (w *WishList) ButtonText() string {
	return w.Name
}

func (w *WishList) Search(ctx context.Context, tx *gorm.DB, params *common.SearchParams) ([]common.PaginatedEntity, error) {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	var wishLists []*WishList
	query := db.Model(&WishList{})

	if params.ListId != "" {
		query = query.Where("user_id = ?", params.ListId)
	}

	if err := query.Find(&wishLists).Error; err != nil {
		slog.Error("Error when searching wish lists: " + err.Error())
		return nil, err
	}

	var entities []common.PaginatedEntity
	for _, wishList := range wishLists {
		entities = append(entities, wishList)
	}

	return entities, nil
}

func (w *WishList) Filter(ctx context.Context, tx *gorm.DB) ([]*WishList, error) {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	var wishLists []*WishList
	query := db.Model(&WishList{})

	if w.UserId != "" {
		query = query.Where("user_id = ?", w.UserId)
	}
	if w.ID != "" {
		query = query.Where("id = ?", w.ID)
	}
	if w.ChatId != "" {
		query = query.Where("chat_id = ?", w.ChatId)
	}

	if err := query.Find(&wishLists).Error; err != nil {
		slog.Error("Error when filtering wish lists: " + err.Error())
		return nil, err
	}

	return wishLists, nil
}

func (w *WishList) Save(ctx context.Context, tx *gorm.DB) error {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	db = db.Session(&gorm.Session{
		SkipHooks: true,
	})

	_, _, _ = w.RefresTimestamps()

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "user_id", "chat_id", "updated_at"}),
	}).Create(w)
	if result.Error != nil {
		slog.Error("Error when trying to save wish list: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("Wish list created/updated")
	return nil
}

func (w *WishList) Delete(ctx context.Context, tx *gorm.DB) error {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	query := db.Model(&WishList{})

	if w.ID != "" {
		query = query.Where("id = ?", w.ID)
	}
	if w.UserId != "" {
		query = query.Where("user_id = ?", w.UserId)
	}
	if w.ChatId != "" {
		query = query.Where("chat_id = ?", w.ChatId)
	}

	result := query.Delete(&WishList{})
	if result.Error != nil {
		slog.Error("Error when trying to delete wish list: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("Wish list deleted")
	return nil
}
