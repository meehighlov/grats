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
		`INSERT INTO user(id, tgid, name, tgusername, chatid, birthday, isadmin, createdat, updatedat)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT(tgid) DO UPDATE SET name=$3, tgusername=$4, chatid=$5, birthday=$6, isadmin=$7, updatedat=$9
        RETURNING id;`,
		&user.ID,
		&user.TGId,
		&user.Name,
		&user.TGusername,
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
	if user.TGId != "" {
		where = append(where, "tgid=$tgid")
	}
	if user.TGusername != "" {
		where = append(where, "tgusername=$tgusername")
	}

	where_ := strings.Join(where, " AND ")
	query := `SELECT id, tgid, name, tgusername, chatid, birthday, isadmin, createdat, updatedat FROM user WHERE ` + where_ + `;`

	rows, err := tx.QueryContext(
		ctx,
		query,
		sql.Named("tgid", user.TGId),
		sql.Named("tgusername", user.TGusername),
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
			&user.TGId,
			&user.Name,
			&user.TGusername,
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

func (friend *Friend) Filter(ctx context.Context, tx *sql.Tx) ([]Friend, error) {
	where := []string{}
	if friend.FilterNotifyAt != "" {
		where = append(where, "notifyat=$notifyat")
	}
	if friend.UserId != "" {
		where = append(where, "userid=$userid")
	}
	if friend.Name != "" {
		where = append(where, "name=$name")
	}
	if friend.ID != "" {
		where = append(where, "id=$id")
	}
	if friend.ChatId != "" {
		where = append(where, "chatid=$chatid")
	}

	where_ := strings.Join(where, " AND ")
	query := `SELECT f.id, f.name, f.birthday, f.userid, f.chatid, f.notifyat, f.createdat, f.updatedat 
              FROM friend f 
              WHERE ` + where_ + `;`

	rows, err := tx.QueryContext(
		ctx,
		query,
		sql.Named("userid", friend.UserId),
		sql.Named("notifyat", friend.FilterNotifyAt),
		sql.Named("name", friend.Name),
		sql.Named("id", friend.ID),
		sql.Named("chatid", friend.ChatId),
	)
	if err != nil {
		slog.Error("Error when filtering friends " + err.Error())
		return nil, err
	}

	friends := []Friend{}

	for rows.Next() {
		friend := Friend{}
		err := rows.Scan(
			&friend.ID,
			&friend.Name,
			&friend.BirthDay,
			&friend.UserId,
			&friend.ChatId,
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
		`INSERT INTO friend(id, name, birthday, userid, chatid, notifyat, createdat, updatedat)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT(id) DO UPDATE SET name=$2, birthday=$3, userid=$4, chatid=$5, notifyat=$6, createdat=$7, updatedat=$8
        RETURNING id;`,
		friend.ID,
		friend.Name,
		friend.BirthDay,
		friend.UserId,
		friend.ChatId,
		*friend.GetNotifyAt(),
		friend.CreatedAt,
		friend.UpdatedAt,
	)
	if err != nil {
		slog.Error("Error when trying to save friend: " + err.Error())
		return err
	}
	slog.Debug("Friend added/updated")

	return nil
}

func (friend *Friend) Delete(ctx context.Context, tx *sql.Tx) error {
	where := []string{}
	if friend.ID != "" {
		where = append(where, "id=$id")
	}
	if friend.ChatId != "" {
		where = append(where, "chatid=$chatid")
	}

	where_ := strings.Join(where, " AND ")
	query := `DELETE FROM friend WHERE ` + where_ + `;`

	_, err := tx.ExecContext(
		ctx,
		query,
		sql.Named("id", friend.ID),
		sql.Named("chatid", friend.ChatId),
	)
	if err != nil {
		slog.Error("Error when trying to delete friend rows: " + err.Error())
		return err
	}

	slog.Debug("Friend row deleted")

	return nil
}

func (friend *Friend) GetTGChatId(ctx context.Context, tx *sql.Tx) (string, error) {
	var tgChatId string
	err := tx.QueryRowContext(
		ctx,
		`SELECT c.tgchatid FROM chat c WHERE c.id = $1`,
		friend.ChatId,
	).Scan(&tgChatId)

	if err != nil {
		slog.Error("Error getting TGChatId for friend: " + err.Error())
		return "", err
	}

	return tgChatId, nil
}

func GetOrCreateChatByTGChatId(ctx context.Context, tx *sql.Tx, tgChatId string, chatType string, botInvitedBy string) (*Chat, error) {
	chat := &Chat{
		TGChatId: tgChatId,
	}

	chats, err := chat.Filter(ctx, tx)
	if err != nil {
		slog.Error("Error getting chat by TGChatId: " + err.Error())
		return nil, err
	}

	if len(chats) > 0 {
		return &chats[0], nil
	}

	chat.BaseFields = NewBaseFields()
	chat.ChatType = chatType
	chat.BotInvitedBy = botInvitedBy
	chat.GreetingTemplate = "🔔Сегодня день рождения у %s🥳"

	err = chat.Save(ctx, tx)
	if err != nil {
		slog.Error("Error creating new chat: " + err.Error())
		return nil, err
	}

	return chat, nil
}

func (c *Chat) Filter(ctx context.Context, tx *sql.Tx) ([]Chat, error) {
	where := []string{}
	if c.ID != "" {
		where = append(where, "id=$id")
	}
	if c.TGChatId != "" {
		where = append(where, "tgchatid=$tgchatid")
	}
	if c.ChatType != "" {
		where = append(where, "chattype LIKE $chattype")
	}
	if c.BotInvitedBy != "" {
		where = append(where, "botinvitedbyid=$botinvitedbyid")
	}

	where_ := strings.Join(where, " AND ")
	query := `SELECT id, tgchatid, chattype, botinvitedbyid, greeting_template, createdat, updatedat FROM chat WHERE ` + where_ + `;`

	rows, err := tx.QueryContext(
		ctx,
		query,
		sql.Named("id", c.ID),
		sql.Named("tgchatid", c.TGChatId),
		sql.Named("chattype", c.ChatType),
		sql.Named("botinvitedbyid", c.BotInvitedBy),
	)
	if err != nil {
		slog.Error("Error when filtering chats " + err.Error())
		return nil, err
	}

	chats := []Chat{}

	for rows.Next() {
		chat := Chat{}
		err := rows.Scan(
			&chat.ID,
			&chat.TGChatId,
			&chat.ChatType,
			&chat.BotInvitedBy,
			&chat.GreetingTemplate,
			&chat.CreatedAt,
			&chat.UpdatedAt,
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
		`INSERT INTO chat(id, tgchatid, chattype, botinvitedbyid, greeting_template, createdat, updatedat)
        VALUES($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT(tgchatid) DO UPDATE SET tgchatid=$2, chattype=$3, botinvitedbyid=$4, greeting_template=$5, updatedat=$7
        RETURNING id;`,
		c.ID,
		c.TGChatId,
		c.ChatType,
		c.BotInvitedBy,
		c.GreetingTemplate,
		c.CreatedAt,
		c.UpdatedAt,
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
		`DELETE FROM chat WHERE tgchatid = $1;`,
		c.TGChatId,
	)
	if err != nil {
		slog.Error("Error when trying to delete chat: " + err.Error())
		return err
	}

	slog.Debug("Chat deleted")

	return nil
}

func (access *Access) All(ctx context.Context, tx *sql.Tx) (*map[string]Access, error) {
	rows, err := tx.QueryContext(ctx, `SELECT id, tgusername FROM access;`)
	if err != nil {
		slog.Error("Error when fetching access list, error: " + err.Error())
		return nil, err
	}
	// defer rows.Close()

	accessList := make(map[string]Access)

	for rows.Next() {
		access := Access{}
		err := rows.Scan(&access.ID, &access.TGusername)
		if err != nil {
			slog.Error("Error when fetching access, error: " + err.Error())
			continue
		}
		accessList[access.TGusername] = access
	}

	return &accessList, nil
}

// idempotent save
// accepts ALL fields of entity and save as is
func (access *Access) Save(ctx context.Context, tx *sql.Tx) error {
	_, _, _ = access.RefresTimestamps()

	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO access(id, tgusername, createdat, updatedat)
        VALUES($1, $2, $3, $4)
        ON CONFLICT(tgusername) DO UPDATE SET tgusername=$2, updatedat=$4
        RETURNING id;`,
		&access.ID,
		&access.TGusername,
		&access.CreatedAt,
		&access.UpdatedAt,
	)
	if err != nil {
		slog.Error("Error when trying to save access: " + err.Error())
		return err
	}
	slog.Debug("Access created/updated")

	return nil
}

func (access *Access) IsExist(ctx context.Context, tx *sql.Tx) bool {
	result := tx.QueryRowContext(
		ctx,
		`SELECT COUNT(id) FROM access WHERE tgusername=$1;`,
		&access.TGusername,
	)
	var count *int
	if err := result.Scan(&count); err != nil {
		slog.Error("Error when trying to get access: " + err.Error())
		return false
	}

	return *count == 1
}

func (access *Access) Delete(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(
		ctx,
		`DELETE FROM access WHERE tgusername = $1;`,
		&access.TGusername,
	)
	if err != nil {
		slog.Error("Error when trying to delete access row: " + err.Error())
		return err
	}

	slog.Debug("Access row deleted")

	return nil
}

// Метод для создания друга с привязкой к чату
func CreateFriendWithChat(ctx context.Context, tx *sql.Tx, name string, birthday string, userId string, tgChatId string) (*Friend, error) {
	// Получаем или создаем чат
	chat, err := GetOrCreateChatByTGChatId(ctx, tx, tgChatId, "private", userId)
	if err != nil {
		slog.Error("Error getting or creating chat: " + err.Error())
		return nil, err
	}

	friend := &Friend{
		BaseFields: NewBaseFields(),
		Name:       name,
		BirthDay:   birthday,
		UserId:     userId,
		ChatId:     chat.ID,
	}

	_, err = friend.RenewNotifayAt()
	if err != nil {
		slog.Error("Error setting notify date: " + err.Error())
		return nil, err
	}

	err = friend.Save(ctx, tx)
	if err != nil {
		slog.Error("Error saving friend: " + err.Error())
		return nil, err
	}

	return friend, nil
}
