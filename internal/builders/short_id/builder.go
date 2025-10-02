package short_id

import (
	"math/rand"
	"time"

	"github.com/meehighlov/grats/internal/config"
)

type Builder struct {
	length  int
	charset string
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func New(cfg *config.Config) *Builder {
	return &Builder{
		length:  cfg.ShortIDLength,
		charset: charset,
	}
}

func (s *Builder) Build() string {
	rand.NewSource(time.Now().UnixNano())
	b := make([]byte, s.length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
