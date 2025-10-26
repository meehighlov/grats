package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	SHORT_ID_LENGTH = 6
)

type BaseFields struct {
	ID        string    `gorm:"primaryKey;type:string;column:id"`
	CreatedAt time.Time `gorm:"not null;column:created_at;type:timestamp with time zone"`
	UpdatedAt time.Time `gorm:"not null;column:updated_at;type:timestamp with time zone"`
}

func (b *BaseFields) RefresTimestamps(tz string) (created time.Time, updated time.Time, _ error) {
	location, err := time.LoadLocation(tz)
	if err != nil {
		return b.CreatedAt, b.UpdatedAt, errors.New("error loading location by timezone, using system timezone, error: " + err.Error() + " entityId: " + b.ID)
	}
	now := time.Now().In(location)
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	b.UpdatedAt = now

	return b.CreatedAt, b.UpdatedAt, nil
}

func NewBaseFields(shortId bool, tz string) (BaseFields, error) {
	id := uuid.New().String()
	if shortId {
		id = GenerateShortID(SHORT_ID_LENGTH)
	}
	location, err := time.LoadLocation(tz)
	if err != nil {
		return BaseFields{}, errors.New("error loading location by timezone, using system timezone, error: " + err.Error() + " NewEntityId: " + id)
	}
	now := time.Now().In(location)
	return BaseFields{id, now, now}, nil
}
