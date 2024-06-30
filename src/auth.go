package src

import (
	"os"
	"strings"
	"github.com/meehighlov/grats/telegram"
)

func IsAuthUser(user telegram.User) bool {
	for _, auth_user_name := range strings.Split(os.Getenv("AUTH_USERS"), ",") {
		if auth_user_name == user.Username {
			return true
		}
	}

	return false
}
