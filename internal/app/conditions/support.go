package conditions

import (
	"context"

	"github.com/meehighlov/grats/pkg/telegram/models"
)

type SupportCondition struct {
	cupportChatId string
}

func SupportReplyCondition(cupportChatId string) *SupportCondition {
	return &SupportCondition{cupportChatId: cupportChatId}
}

func (c *SupportCondition) Check(ctx context.Context, update *models.Update) (bool, error) {
	if update.GetMessage() != nil &&
		update.GetMessage().GetChatIdStr() == c.cupportChatId &&
		update.GetMessage().IsReply() {
		return true, nil
	}

	return false, nil
}
