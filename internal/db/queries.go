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

func (friend *Friend) GetId() string {
	return friend.ID
}

func (chat *Chat) GetId() string {
	return chat.ID
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
		DoUpdates: clause.AssignmentColumns([]string{"name", "tg_username", "chat_id", "birthday", "is_admin", "updated_at"}),
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

func (friend *Friend) Filter(ctx context.Context, tx *gorm.DB) ([]*Friend, error) {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	var friends []*Friend
	query := db.Model(&Friend{})

	if friend.FilterNotifyAt != "" {
		query = query.Where("notify_at = ?", friend.FilterNotifyAt)
	}
	if friend.UserId != "" {
		query = query.Where("user_id = ?", friend.UserId)
	}
	if friend.Name != "" {
		query = query.Where("name = ?", friend.Name)
	}
	if friend.ID != "" {
		query = query.Where("id = ?", friend.ID)
	}
	if friend.ChatId != "" {
		query = query.Where("chat_id = ?", friend.ChatId)
	}

	if err := query.Find(&friends).Error; err != nil {
		slog.Error("Error when filtering friends: " + err.Error())
		return nil, err
	}

	return friends, nil
}

func (friend *Friend) Save(ctx context.Context, tx *gorm.DB) error {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	db = db.Session(&gorm.Session{
		SkipHooks: true,
	})

	_, _, _ = friend.RefresTimestamps()

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "user_id", "birthday", "chat_id", "notify_at", "updated_at"}),
	}).Create(friend)
	if result.Error != nil {
		slog.Error("Error when trying to save friend: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("Friend created/updated")
	return nil
}

func (friend *Friend) Delete(ctx context.Context, tx *gorm.DB) error {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	query := db.Model(&Friend{})

	if friend.ID != "" {
		query = query.Where("id = ?", friend.ID)
	}
	if friend.UserId != "" {
		query = query.Where("user_id = ?", friend.UserId)
	}
	if friend.ChatId != "" {
		query = query.Where("chat_id = ?", friend.ChatId)
	}

	result := query.Delete(&Friend{})
	if result.Error != nil {
		slog.Error("Error when trying to delete friend: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("Friend deleted")
	return nil
}

func (f *Friend) GreaterThan(other common.PaginatedEntity) bool {
	otherFriend, ok := other.(*Friend)
	if !ok {
		return false
	}

	if f.IsTodayBirthday() {
		return true
	}
	if otherFriend.IsTodayBirthday() {
		return false
	}
	countI := f.CountDaysToBirthday()
	countJ := otherFriend.CountDaysToBirthday()
	return countI > countJ
}

func (c *Chat) Filter(ctx context.Context, tx *gorm.DB) ([]*Chat, error) {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	var chats []*Chat
	query := db.Model(&Chat{})

	if c.ChatId != "" {
		query = query.Where("chat_id = ?", c.ChatId)
	}
	if c.ID != "" {
		query = query.Where("id = ?", c.ID)
	}
	if c.BotInvitedById != "" {
		query = query.Where("bot_invited_by_id = ?", c.BotInvitedById)
	}

	if err := query.Find(&chats).Error; err != nil {
		slog.Error("Error when filtering chats: " + err.Error())
		return nil, err
	}

	return chats, nil
}

func (c *Chat) Save(ctx context.Context, tx *gorm.DB) error {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	db = db.Session(&gorm.Session{
		SkipHooks: true,
	})

	_, _, _ = c.RefresTimestamps()

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chat_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"chat_type", "bot_invited_by_id", "greeting_template", "silent_notifications", "updated_at"}),
	}).Create(c)
	if result.Error != nil {
		slog.Error("Error when trying to save chat: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("Chat created/updated")
	return nil
}

func (c *Chat) Delete(ctx context.Context, tx *gorm.DB) error {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	query := db.Model(&Chat{})

	if c.ID != "" {
		query = query.Where("id = ?", c.ID)
	}
	if c.ChatId != "" {
		query = query.Where("chat_id = ?", c.ChatId)
	}
	if c.BotInvitedById != "" {
		query = query.Where("bot_invited_by_id = ?", c.BotInvitedById)
	}

	result := query.Delete(&Chat{})
	if result.Error != nil {
		slog.Error("Error when trying to delete chat: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("Chat deleted")
	return nil
}

func (c *Chat) GreaterThan(other common.PaginatedEntity) bool {
	return true
}

func (f *Friend) ButtonText() string {
	buttonText := fmt.Sprintf("%s %s", f.Name, f.BirthDay)

	if f.IsTodayBirthday() {
		buttonText = fmt.Sprintf("%s 🥳", buttonText)
	} else {
		if f.IsThisMonthAfterToday() {
			buttonText = fmt.Sprintf("%s 🕒", buttonText)
		}
	}

	return buttonText
}

func (c *Chat) ButtonText() string {
	return c.ChatId
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

func (f *Friend) Search(ctx context.Context, tx *gorm.DB, params *common.SearchParams) ([]common.PaginatedEntity, error) {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	var friends []*Friend
	query := db.Model(&Friend{})

	if params.ListId != "" {
		query = query.Where("chat_id = ?", params.ListId)
	}

	if err := query.Find(&friends).Error; err != nil {
		slog.Error("Error when searching friends: " + err.Error())
		return nil, err
	}

	var entities []common.PaginatedEntity
	for _, friend := range friends {
		entities = append(entities, friend)
	}

	return entities, nil
}

func (c *Chat) Search(ctx context.Context, tx *gorm.DB, params *common.SearchParams) ([]common.PaginatedEntity, error) {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	var chats []*Chat
	query := db.Model(&Chat{})

	if params.ListId != "" {
		query = query.Where("chat_id = ?", params.ListId)
	}

	if err := query.Find(&chats).Error; err != nil {
		slog.Error("Error when searching chats: " + err.Error())
		return nil, err
	}

	var entities []common.PaginatedEntity
	for _, chat := range chats {
		entities = append(entities, chat)
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
	if w.ExecutorId == "" && otherWish.ExecutorId != "" {
		return true
	}
	if w.ExecutorId != "" && otherWish.ExecutorId == "" {
		return false
	}
	if w.Link != "" || otherWish.Link != "" {
		return false
	}
	// TODO make it with int
	if w.Price != "" && otherWish.Price != "" {
		return w.Price < otherWish.Price
	}
	if w.Price != "" || otherWish.Price != "" {
		return false
	}
	return false
}

func (w *Wish) ButtonText() string {
	text := w.Name
	if w.Price != "" {
		text = w.Price + "(RUB)" + " " + w.Name
	}
	if w.ExecutorId != "" {
		text = "🔒 " + text
	}
	return text
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
	query := db.Model(&Wish{}).Clauses(clause.Locking{Strength: "UPDATE"})

	if w.ID != "" {
		query = query.Where("id = ?", w.ID)
	}

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
		DoUpdates: clause.AssignmentColumns([]string{"name", "chat_id", "user_id", "wish_list_id", "link", "price", "executor_id", "updated_at"}),
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

	db = db.Session(&gorm.Session{
		SkipHooks: true,
	})

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
		query = query.Where("id = ?", params.ListId)
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
		slog.Error("Error when trying to save wishList: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("WishList created/updated")
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
	if w.Name != "" {
		query = query.Where("name = ?", w.Name)
	}

	result := query.Delete(&WishList{})
	if result.Error != nil {
		slog.Error("Error when trying to delete wishList: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("WishList deleted")
	return nil
}
