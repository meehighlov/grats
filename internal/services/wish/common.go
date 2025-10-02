package wish

import (
	"context"
	"errors"

	inlinekeyboard "github.com/meehighlov/grats/internal/builders/inline_keyboard"
	"github.com/meehighlov/grats/internal/repositories/entities"
)

func (s *Service) PickFirstWishList(ctx context.Context, userId string) (*entities.WishList, error) {
	filter := entities.WishList{UserId: userId}
	wishList, err := s.repositories.WishList.Filter(ctx, nil, &filter)
	if err != nil {
		return nil, err
	}

	if len(wishList) == 0 {
		return nil, errors.New("wish list is empty")
	}

	return wishList[0], nil
}

func (s *Service) BuildEntityButtons(wishes []*entities.Wish, offset int, callback func(id string, offset int) string) *inlinekeyboard.Builder {
	buttons := s.builders.KeyboardBuilder.NewKeyboard()
	for _, entity := range wishes {
		buttonText := entity.ButtonText()

		buttons.AppendAsLine(buttons.NewButton(buttonText, callback(entity.ID, offset)))
	}

	return buttons
}
