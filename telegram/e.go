package telegram

import (
	"fmt"
)

func wrapErr(msg string, err error) error {
	err_ := fmt.Errorf("%s: %w", msg, err)
	return err_
}

func wrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	return wrapErr(msg, err)
}
