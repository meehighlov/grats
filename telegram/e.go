package telegram

import "fmt"

func wrapErr(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func wrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	return wrapErr(msg, err)
}
