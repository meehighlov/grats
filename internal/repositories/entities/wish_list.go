package entities

type WishList struct {
	BaseFields

	Name   string `gorm:"column:name;type:varchar"`
	UserId string `gorm:"not null;index;column:user_id;type:varchar"`
	ChatId string `gorm:"column:chat_id;type:varchar"`

	User User `gorm:"foreignKey:UserId;references:ID"`
}

func (WishList) TableName() string {
	return "wish_list"
}

func (wishList *WishList) GetUserId() string {
	return wishList.UserId
}
