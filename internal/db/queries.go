package db

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

func (user *User) Filter(ctx context.Context, tx *gorm.DB) ([]User, error) {
	db := GetDB()
	if tx != nil {
		db = tx
	} else {
		db = db.WithContext(ctx)
	}

	var users []User
	query := db.Model(&User{})

	if user.TgId != "" {
		query = query.Where("tg_id = ?", user.TgId)
	}
	if user.TgUsername != "" {
		query = query.Where("tg_username = ?", user.TgUsername)
	}

	if err := query.Find(&users).Error; err != nil {
		slog.Error("Error when filtering users: " + err.Error())
		return nil, err
	}

	return users, nil
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

	result := db.Delete(friend)
	if result.Error != nil {
		slog.Error("Error when trying to delete friend: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("Friend deleted")
	return nil
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

	result := db.Delete(c)
	if result.Error != nil {
		slog.Error("Error when trying to delete chat: " + result.Error.Error())
		return result.Error
	}

	slog.Debug("Chat deleted")
	return nil
}
