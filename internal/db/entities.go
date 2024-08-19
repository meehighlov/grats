package db

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/meehighlov/grats/internal/config"
)

type BaseFields struct {
	ID        string
	CreatedAt string
	UpdatedAt string
}

func (b *BaseFields) RefresTimestamps() (created string, updated string, _ error) {
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " entityId: " + b.ID)
	}
	now := time.Now().In(location).Format("02.01.2006T15:04:05")
	if b.CreatedAt == "" {
		b.CreatedAt = now
	}
	b.UpdatedAt = now

	return b.CreatedAt, b.UpdatedAt, nil
}

func NewBaseFields() BaseFields {
	id := uuid.New().String()
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " NewEntityId: " + id)
	}
	now := time.Now().In(location).Format("02.01.2006T15:04:05")
	return BaseFields{id, now, now}
}

type User struct {
	// telegram user -> bot's user

	BaseFields

	TGId       int // id will be taken from telegram
	Name       string
	TGusername string
	ChatId     int // chatId - id of chat with user, bot uses it to send notification
	Birthday   string
	IsAdmin    int

	Friends []Friend
}

func (user *User) HasAdminAccess() bool {
	return user.IsAdmin == 1
}

func (user *User) GetTGUserName() string {
	if !strings.HasPrefix("@", user.TGusername) {
		return "@" + user.TGusername
	}
	return user.TGusername
}

func (user *User) FriendsListAsString() string {
	result := ""
	for _, friend := range user.Friends {
		result += friend.Name + " " + friend.BirthDay + "\n"
	}
	return result
}

type Friend struct {
	BaseFields

	// todo store timezone in friend table or somewere in db - for user's specific timezone

	Name           string
	UserId         int
	BirthDay       string
	ChatId         int
	notifyAt       string
	FilterNotifyAt string // this param is only for filtering
}

func (friend *Friend) CountDaysToBirthday() int {
	// todo store timezone in friend table or somewere in db - for user's specific timezone
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " friendId: " + friend.ID)
	}
	now := time.Now().In(location)
	notify, err := time.Parse("02.01.2006", *friend.GetNotifyAt())
	if err != nil {
		slog.Error("error parsing notify during count days to birthday: " + err.Error())
		return -1
	}

	diff := now.Sub(notify)
	diff_days := diff.Hours() / 24

	return int(diff_days)
}

func (friend *Friend) IsThisMonthAfterToday() bool {
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " friendId: " + friend.ID)
	}
	now := strings.Split(time.Now().In(location).Format("02.01.2006"), ".")
	thisMonth := strings.Split(friend.BirthDay, ".")[1] == now[1]
	afterToday := strings.Split(friend.BirthDay, ".")[0] > now[0]

	return thisMonth && afterToday
}

func (friend *Friend) IsTodayBirthday() bool {
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " friendId: " + friend.ID)
	}
	now := strings.Split(time.Now().In(location).Format("02.01.2006"), ".")
	bd := strings.Split(friend.BirthDay, ".")

	return now[0] == bd[0] && now[1] == bd[1]
}

func (friend *Friend) GetNotifyAt() *string {
	return &friend.notifyAt
}

func (friend *Friend) RenewNotifayAt() (string, error) {
	format := "02.01" // day.month

	parts := strings.Split(friend.BirthDay, ".")
	birtday_wo_year := strings.Join(parts[:2], ".")

	birthday, err := time.Parse(format, birtday_wo_year)

	if err != nil {
		slog.Error("notify date creation: cannot parse birthday:" + err.Error())
		return "", nil
	}

	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " friendId: " + friend.ID)
	}

	today, err := time.Parse(format, time.Now().In(location).Format(format))

	if err != nil {
		slog.Error("notify date creation: cannot parse today date:" + err.Error())
		return "", nil
	}

	year := time.Now().In(location).Year()
	if today.After(birthday) || today.Equal(birthday) {
		year += 1
	}

	*friend.GetNotifyAt() = fmt.Sprintf(birthday.Format(format)+".%d", year)

	return *friend.GetNotifyAt(), nil
}

func (friend *Friend) UpdateNotifyAt() (string, error) {
	format := "02.01.2006" // day.month.year
	notifyAt, err := time.Parse(format, *friend.GetNotifyAt())

	if err != nil {
		slog.Error("error updating notfiy date: " + err.Error())
		return *friend.GetNotifyAt(), err
	}

	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " friendId: " + friend.ID)
	}

	// in case of duplicated call
	if notifyAt.Year() != time.Now().In(location).Year() {
		return *friend.GetNotifyAt(), nil
	}

	*friend.GetNotifyAt() = notifyAt.AddDate(1, 0, 0).Format(format)
	return *friend.GetNotifyAt(), nil
}

func (friend *Friend) GetChatIdStr() string {
	return strconv.Itoa(friend.ChatId)
}

type Access struct {
	BaseFields

	TGusername string
}

func (access *Access) GetTGUserName() string {
	if !strings.HasPrefix("@", access.TGusername) {
		return "@" + access.TGusername
	}
	return access.TGusername
}
