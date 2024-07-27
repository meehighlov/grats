package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type BaseFields struct {
	ID string
	CreatedAt string
	UpdatedAt string
}

func (b *BaseFields) RefresTimestamps() (created string, updated string, _ error) {
	now := time.Now().Format("02.01.2006T15:04:05")
	if b.CreatedAt == "" {
		b.CreatedAt = now
	}
	b.UpdatedAt = now

	return b.CreatedAt, b.UpdatedAt, nil
}

func NewBaseFields() BaseFields {
	now := time.Now().Format("02.01.2006T15:04:05")
	return BaseFields{uuid.New().String(), now, now}
}

type User struct {
	// telegram user -> bot's user

	BaseFields

	TGId       int  // id will be taken from telegram
	Name       string
	TGusername string
	ChatId     int  // chatId - id of chat with user, bot uses it to send notification
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

	Name     string
	UserId   int
	BirthDay string
	ChatId   int
	notifyAt string
	FilterNotifyAt string  // this params is only for filtering
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
		log.Println("notify date creation: cannot parse birthday:", err.Error())
		return "", nil
	}

	today, err := time.Parse(format, time.Now().Format(format))

	if err != nil {
		log.Println("notify date creation: cannot parse today date:", err.Error())
		return "", nil
	}

	year := time.Now().Year()
	if today.After(birthday) || today.Equal(birthday) {
		year += 1
	}

	*friend.GetNotifyAt() = fmt.Sprintf(birthday.Format(format)+".%d", year)

	return *friend.GetNotifyAt(), nil
}

func (friend *Friend) NotifyNeeded() (bool, error) {
	format := "02.01.2006" // day.month.year
	today, err := time.Parse(format, time.Now().Format(format))

	if err != nil {
		log.Println("notify check error: cannot parse today date:", err.Error())
		return false, err
	}

	notifyAt, err := time.Parse(format, *friend.GetNotifyAt())

	if err != nil {
		log.Println("notify check error: cannot parse notify date:", err.Error())
		return false, err
	}

	if today.Equal(notifyAt) {
		return true, nil
	}

	return false, nil
}

func (friend *Friend) UpdateNotifyAt() (string, error) {
	format := "02.01.2006" // day.month.year
	notifyAt, err := time.Parse(format, *friend.GetNotifyAt())

	if err != nil {
		log.Println("error updating notfiy date:", err.Error())
		return *friend.GetNotifyAt(), err
	}

	// in case of duplicated call
	if notifyAt.Year() != time.Now().Year() {
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
