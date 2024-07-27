package db

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
)

// idempotent save
// accepts ALL fields of entity and save as is
func (user *User) Save(ctx context.Context) error {
	_, _, _ = user.RefresTimestamps()

	_, err := sqliteConn.ExecContext(
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

func (user *User) Filter(ctx context.Context) ([]User, error) {
	where := []string{}
	if user.TGId != 0 {
		where = append(where, "tgid=$tgid")
	}
	if user.TGusername != "" {
		where = append(where, "tgusername=$tgusername")
	}

	where_ := strings.Join(where, " AND ")
	query := `SELECT id, tgid, name, tgusername, chatid, birthday, isadmin, createdat, updatedat FROM user WHERE ` + where_ + `;`

	rows, err := sqliteConn.QueryContext(
		ctx,
		query,
		sql.Named("tgid", user.TGId),
		sql.Named("tgusername", user.TGusername),
	)
	if err != nil {
		slog.Error("Error when filtering users " + err.Error())
		return nil, err
	}
	defer rows.Close()

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

func (friend *Friend) Filter(ctx context.Context) ([]Friend, error) {
	where := []string{}
	if friend.FilterNotifyAt != "" {
		where = append(where, "notifyat=$notifyat")
	}
	if friend.UserId != 0 {
		where = append(where, "userid=$userid")
	}

	where_ := strings.Join(where, " AND ")
	query := `SELECT id, name, birthday, userid, chatid, notifyat, createdat, updatedat FROM friend WHERE ` + where_ + `;`

	rows, err := sqliteConn.QueryContext(
		ctx,
		query,
		sql.Named("userid", friend.UserId),
		sql.Named("notifyat", friend.FilterNotifyAt),
	)
	if err != nil {
		slog.Error("Error when filtering friends " + err.Error())
		return nil, err
	}
	defer rows.Close()

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
func (friend *Friend) Save(ctx context.Context) error {
	_, _, _ = friend.RefresTimestamps()

	_, err := sqliteConn.ExecContext(
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

func (access *Access) All(ctx context.Context) (*map[string]Access, error) {
	rows, err := sqliteConn.QueryContext(ctx, `SELECT id, tgusername FROM access;`)
	if err != nil {
		slog.Error("Error when fetching access list, error: " + err.Error())
		return nil, err
	}
	defer rows.Close()

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
func (access *Access) Save(ctx context.Context) error {
	_, _, _ = access.RefresTimestamps()

	_, err := sqliteConn.ExecContext(
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

func (access *Access) IsExist(ctx context.Context) bool {
	result := sqliteConn.QueryRowContext(
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

func (access *Access) Delete(ctx context.Context) error {
	_, err := sqliteConn.ExecContext(
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
