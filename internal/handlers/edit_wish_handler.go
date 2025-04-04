package handlers

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

const (
	WISH_LINK_MAX_LEN = 500
)

func EditPriceHandler(ctx context.Context, event *common.Event) error {
	event.ReplyCallbackQuery(ctx, "Введите цену в рублях")

	event.GetContext().AppendText(event.GetCallbackQuery().Data)

	refreshMessageId := event.GetCallbackQuery().Message.GetMessageIdStr()
	event.GetContext().AppendText(refreshMessageId)

	event.SetNextHandler("edit_price_save")

	return nil
}

func SaveEditPriceHandler(ctx context.Context, event *common.Event) error {
	message := event.GetMessage()
	params := common.CallbackFromString(event.GetContext().GetTexts()[0])
	wishId := params.Id

	done := false

	err := db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		wish, err := (&db.Wish{BaseFields: db.BaseFields{ID: wishId}}).Filter(ctx, tx)
		if err != nil {
			event.Reply(ctx, "Не смог найти желание")
			return err
		}

		price, err := strconv.ParseFloat(message.Text, 64)
		if err != nil || price <= 0 {
			event.Reply(ctx, "Не могу распознать цену, пример 👉 1000, 400.5")
			event.SetNextHandler("edit_price_save")
			return nil
		}

		wish[0].Price = message.Text
		if err := wish[0].Save(ctx, tx); err != nil {
			return err
		}

		done = true

		executorId := strconv.Itoa(event.GetMessage().From.Id)

		chatId := event.GetMessage().GetChatIdStr()
		refreshMessageId := event.GetContext().GetTexts()[1]
		event.RefreshMessage(ctx, chatId, refreshMessageId, wish[0].Info(executorId), *buildWishInfoKeyboard(wish[0], fmt.Sprintf("%d", LIST_START_OFFSET), params.Pagination.Direction, params.SourceId).Murkup())

		return nil
	})

	if err != nil {
		return err
	}

	if done {
		event.Reply(ctx, "Цена установлена 💾")
		event.SetNextHandler("")
	}

	return nil
}

func EditLinkHandler(ctx context.Context, event *common.Event) error {
	event.ReplyCallbackQuery(ctx, "Введите ссылку")

	event.GetContext().AppendText(event.GetCallbackQuery().Data)

	refreshMessageId := event.GetCallbackQuery().Message.GetMessageIdStr()
	event.GetContext().AppendText(refreshMessageId)

	event.SetNextHandler("edit_link_save")

	return nil
}

func SaveEditLinkHandler(ctx context.Context, event *common.Event) error {
	message := event.GetMessage()
	params := common.CallbackFromString(event.GetContext().GetTexts()[0])
	wishId := params.Id

	done := false

	link := message.Text
	if len(link) > WISH_LINK_MAX_LEN {
		event.Reply(ctx, fmt.Sprintf("Ссылка не может быть длиннее %d символов", WISH_LINK_MAX_LEN))
		event.SetNextHandler("edit_link_save")
		return nil
	}

	parsedURL, err := url.Parse(link)
	if err != nil || parsedURL.Host == "" {
		event.Reply(ctx, "Не могу распознать ссылку, введите ссылку в формате https://example.com")
		event.SetNextHandler("edit_link_save")
		return nil
	}

	info, err := getCertificateInfo(link)
	if err != nil {
		event.Reply(ctx, "Я не доверяю этому сайту😔 попробуй другую ссылку")
		event.SetNextHandler("edit_link_save")
		return nil
	}

	event.Logger.Debug("certificate check", "info", info)

	err = db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		wish, err := (&db.Wish{BaseFields: db.BaseFields{ID: wishId}}).Filter(ctx, tx)
		if len(wish) == 0 || err != nil {
			event.Reply(ctx, "Не смог найти желание, либо оно было исполнено")
			return err
		}
		wish[0].Link = link
		if err := wish[0].Save(ctx, tx); err != nil {
			return err
		}

		done = true

		executorId := strconv.Itoa(event.GetMessage().From.Id)
		refreshMessageId := event.GetContext().GetTexts()[1]
		chatId := event.GetMessage().GetChatIdStr()
		event.RefreshMessage(ctx, chatId, refreshMessageId, wish[0].Info(executorId), *buildWishInfoKeyboard(wish[0], fmt.Sprintf("%d", LIST_START_OFFSET), params.Pagination.Direction, params.SourceId).Murkup())

		// delete message with link - it's too large
		event.DeleteMessage(ctx, chatId, message.GetMessageIdStr())

		return nil
	})

	if err != nil {
		return err
	}

	if done {
		event.Reply(ctx, "Ссылка установлена 💾")
		event.SetNextHandler("")
	}

	return nil
}

func EditWishNameHandler(ctx context.Context, event *common.Event) error {
	event.ReplyCallbackQuery(ctx, "Введите новое название желания")

	event.GetContext().AppendText(event.GetCallbackQuery().Data)

	refreshMessageId := event.GetCallbackQuery().Message.GetMessageIdStr()
	event.GetContext().AppendText(refreshMessageId)

	event.SetNextHandler("edit_wish_name_save")

	return nil
}

func SaveEditWishNameHandler(ctx context.Context, event *common.Event) error {
	message := event.GetMessage()
	params := common.CallbackFromString(event.GetContext().GetTexts()[0])
	wishId := params.Id

	done := false

	_, err := validateWishName(message.Text)
	if err != nil {
		event.Reply(ctx, "Имя желания содержит недопустимые символы, попробуйте использовать только цифры и буквы")
		event.SetNextHandler("edit_wish_name_save")
		return nil
	}

	err = db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		wish, err := (&db.Wish{BaseFields: db.BaseFields{ID: wishId}}).Filter(ctx, tx)
		if err != nil {
			event.Logger.Error(
				"SaveEditWishNameHandler",
				"details", err.Error(),
				"chatId", event.GetMessage().GetChatIdStr(),
			)
			event.Reply(ctx, "Не смог найти желание, либо оно было исполнено")
			return err
		}

		wish[0].Name = message.Text
		err = wish[0].Save(ctx, tx)
		if err != nil {
			return err
		}

		done = true

		executorId := strconv.Itoa(event.GetMessage().From.Id)

		refreshMessageId := event.GetContext().GetTexts()[1]
		chatId := event.GetMessage().GetChatIdStr()
		event.RefreshMessage(ctx, chatId, refreshMessageId, wish[0].Info(executorId), *buildWishInfoKeyboard(wish[0], fmt.Sprintf("%d", LIST_START_OFFSET), params.Pagination.Direction, params.SourceId).Murkup())

		return nil
	})

	if err != nil {
		return err
	}

	if done {
		event.Reply(ctx, "Название желания изменено 💾")
		event.SetNextHandler("")
	}

	return nil
}

func getCertificateInfo(urlStr string) (string, error) {
	if !strings.HasPrefix(urlStr, "https://") {
		return "", fmt.Errorf("URL должен начинаться с https://")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	host := parsedURL.Host
	if !strings.Contains(host, ":") {
		host = host + ":443"
	}

	conn, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		return "", fmt.Errorf("ошибка соединения TLS: %v", err)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return "", fmt.Errorf("сертификаты не найдены")
	}

	cert := certs[0]

	certInfo := fmt.Sprintf(
		"Субъект: %s\nИздатель: %s\nДействителен с: %s\nДействителен до: %s\nSAN: %v",
		cert.Subject,
		cert.Issuer,
		cert.NotBefore.Format("2006-01-02"),
		cert.NotAfter.Format("2006-01-02"),
		cert.DNSNames,
	)

	return certInfo, nil
}
