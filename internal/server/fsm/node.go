package fsm

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

type handler func(ctx context.Context, update *telegram.Update) error

type HandledError struct {
	handlerError error
	callNext     string
}

type nodeOption func(*node) error

type node struct {
	command string
	handler handler

	errors  []*HandledError

	nextHandlerName string
}

func new(command string, handler handler) *node {
	return &node{
		command: command,
		handler: handler,
	}
}

func (n *node) getNextHandlerByError(err error) string {
	handler := ""
	for _, e := range n.errors {
		if e.handlerError == err {
			return e.callNext
		}
	}
	return handler
}

func NextHandler(handlerName string) nodeOption {
	return func(n *node) error {
		n.nextHandlerName = handlerName
		return nil
	}
}

func OnErrorCallNext(err error, nextHandlerName string) nodeOption {
	return func(n *node) error {
		handled := &HandledError{
			handlerError: err,
			callNext: nextHandlerName,
		}
		n.errors = append(n.errors, handled)
		return nil
	}
}
