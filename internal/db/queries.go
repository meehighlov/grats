package db

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
)

// idempotent save
// accepts ALL fields of entity and save as is
func (user *User) Save(ctx context.Context, tx *sql.Tx) error {
	_, _, _ = user.RefresTimestamps()

	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO user(id, tg_id, name, tg_username, chat_id, birthday, is_admin, created_at, updated_at)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT(tg_id) DO UPDATE SET name=$3, tg_username=$4, chat_id=$5, birthday=$6, is_admin=$7, updated_at=$9
        RETURNING id;`,
		&user.ID,
		&user.TgId,
		&user.Name,
		&user.TgUsername,
		&user.ChatId,
		&user.Birthday,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		slog.Error("Error when trying to save user: " + err.Error())
		return err
	}
	slog.Debug("User created/updated")

	return nil
}

func (user *User) Filter(ctx context.Context, tx *sql.Tx) ([]User, error) {
	where := []string{}
	if user.TgId != "" {
		where = append(where, "tg_id=$tg_id")
	}
	if user.TgUsername != "" {
		where = append(where, "tg_username=$tg_username")
	}

	where_ := strings.Join(where, " AND ")
	query := `SELECT id, tg_id, name, tg_username, chat_id, birthday, is_admin, created_at, updated_at FROM user WHERE ` + where_ + `;`

	rows, err := tx.QueryContext(
		ctx,
		query,
		sql.Named("tg_id", user.TgId),
		sql.Named("tg_username", user.TgUsername),
	)
	if err != nil {
		slog.Error("Error when filtering users " + err.Error())
		return nil, err
	}

	users := []User{}

	for rows.Next() {
		user := User{}
		err := rows.Scan(
			&user.ID,
			&user.TgId,
			&user.Name,
			&user.TgUsername,
			&user.ChatId,
			&user.Birthday,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			slog.Error("Error fetching users by filter params: " + err.Error())
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func (friend *Friend) Filter(ctx context.Context, tx *sql.Tx) ([]*Friend, error) {
	where := []string{}
	if friend.FilterNotifyAt != "" {
		where = append(where, "notify_at=$notify_at")
	}
	if friend.UserId != "" {
		where = append(where, "user_id=$user_id")
	}
	if friend.Name != "" {
		where = append(where, "name=$name")
	}
	if friend.ID != "" {
		where = append(where, "id=$id")
	}

	where_ := strings.Join(where, " AND ")
	query := `SELECT id, name, birthday, chat_id, user_id, notify_at, created_at, updated_at FROM friend WHERE ` + where_ + `;`

	rows, err := tx.QueryContext(
		ctx,
		query,
		sql.Named("notify_at", friend.FilterNotifyAt),
		sql.Named("user_id", friend.UserId),
		sql.Named("name", friend.Name),
		sql.Named("id", friend.ID),
		sql.Named("chat_id", friend.ChatId),
	)
	if err != nil {
		slog.Error("Error when filtering friends " + err.Error())
		return nil, err
	}

	friends := []*Friend{}

	for rows.Next() {
		friend := &Friend{}
		err := rows.Scan(
			&friend.ID,
			&friend.Name,
			&friend.BirthDay,
			&friend.ChatId,
			&friend.UserId,
			friend.GetNotifyAt(),
			&friend.CreatedAt,
			&friend.UpdatedAt,
		)
		if err != nil {
			slog.Error("Error fetching friends by filter params: " + err.Error())
			continue
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

// idempotent save
// accepts ALL fields of entity and save as is
func (friend *Friend) Save(ctx context.Context, tx *sql.Tx) error {
	_, _, _ = friend.RefresTimestamps()

	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO friend(id, name, birthday, chat_id, user_id, notify_at, created_at, updated_at)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT(name,chat_id) DO UPDATE SET birthday=$3, notify_at=$6, updated_at=$8
        RETURNING id;`,
		friend.ID,
		friend.Name,
		friend.BirthDay,
		friend.ChatId,
		friend.UserId,
		*friend.GetNotifyAt(),
		friend.CreatedAt,
		friend.UpdatedAt,
	)
	if err != nil {
		slog.Error("Error when trying to save friend: " + err.Error())
		return err
	}
	slog.Debug("Friend created/updated")

	return nil
}

func (friend *Friend) Delete(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(
		ctx,
		`DELETE FROM friend WHERE id=$1;`,
		&friend.ID,
	)
	if err != nil {
		slog.Error("Error when trying to delete friend: " + err.Error())
		return err
	}
	slog.Debug("Friend deleted")

	return nil
}

func (c *Chat) Filter(ctx context.Context, tx *sql.Tx) ([]*Chat, error) {
	where := []string{}
	if c.ChatId != "" {
		where = append(where, "chat_id=$chat_id")
	}
	if c.ID != "" {
		where = append(where, "id=$id")
	}
	if c.BotInvitedById != "" {
		where = append(where, "bot_invited_by_id=$bot_invited_by_id")
	}

	where_ := ""
	if len(where) > 0 {
		where_ = "WHERE " + strings.Join(where, " AND ")
	}
	query := `SELECT id, chat_id, chat_type, bot_invited_by_id, created_at, updated_at, greeting_template, silent_notifications FROM chat ` + where_ + `;`

	rows, err := tx.QueryContext(
		ctx,
		query,
		sql.Named("chat_id", c.ChatId),
		sql.Named("id", c.ID),
		sql.Named("bot_invited_by_id", c.BotInvitedById),
	)
	if err != nil {
		slog.Error("Error when filtering chats " + err.Error())
		return nil, err
	}

	chats := []*Chat{}

	for rows.Next() {
		chat := &Chat{}
		err := rows.Scan(
			&chat.ID,
			&chat.ChatId,
			&chat.ChatType,
			&chat.BotInvitedById,
			&chat.CreatedAt,
			&chat.UpdatedAt,
			&chat.GreetingTemplate,
			&chat.SilentNotifications,
		)
		if err != nil {
			slog.Error("Error fetching chats by filter params: " + err.Error())
			continue
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

// idempotent save
// accepts ALL fields of entity and save as is
func (c *Chat) Save(ctx context.Context, tx *sql.Tx) error {
	_, _, _ = c.RefresTimestamps()

	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO chat(id, chat_id, chat_type, bot_invited_by_id, created_at, updated_at, greeting_template, silent_notifications)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT(chat_id) DO UPDATE SET chat_type=$3, bot_invited_by_id=$4, updated_at=$6, greeting_template=$7, silent_notifications=$8
        RETURNING id;`,
		&c.ID,
		&c.ChatId,
		&c.ChatType,
		&c.BotInvitedById,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.GreetingTemplate,
		&c.SilentNotifications,
	)
	if err != nil {
		slog.Error("Error when trying to save chat: " + err.Error())
		return err
	}
	slog.Debug("Chat created/updated")

	return nil
}

func (c *Chat) Delete(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(
		ctx,
		`DELETE FROM chat WHERE id=$1;`,
		&c.ID,
	)
	if err != nil {
		slog.Error("Error when trying to delete chat: " + err.Error())
		return err
	}
	slog.Debug("Chat deleted")

	return nil
}
