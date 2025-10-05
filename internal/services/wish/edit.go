package wish

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/entities"
)

func (s *Service) EditPriceHandler(ctx context.Context, update *telegram.Update) error {
	s.clients.Telegram.Reply(ctx, s.constants.ENTER_PRICE, update)

	s.clients.Cache.AppendText(ctx, update.GetChatIdStr(), update.CallbackQuery.Data)

	refreshMessageId := update.CallbackQuery.Message.GetMessageIdStr()
	s.clients.Cache.AppendText(ctx, update.GetChatIdStr(), refreshMessageId)

	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), s.constants.CMD_EDIT_PRICE_SAVE)

	return nil
}

func (s *Service) SaveEditPriceHandler(ctx context.Context, update *telegram.Update) error {
	message := update.GetMessage()
	texts, err := s.clients.Cache.GetTexts(ctx, update.GetChatIdStr())
	if err != nil {
		return err
	}
	params := s.builders.CallbackDataBuilder.FromString(texts[0])
	wishId := params.ID

	wish, err := s.repositories.Wish.Filter(ctx, &entities.Wish{BaseFields: entities.BaseFields{ID: wishId}})
	if err != nil {
		s.clients.Telegram.Reply(ctx, s.constants.WISH_NOT_FOUND, update)
		return err
	}

	price, err := strconv.ParseFloat(message.Text, 64)
	if err != nil || price <= 0 {
		s.clients.Telegram.Reply(ctx, s.constants.PRICE_INVALID_FORMAT, update)
		s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), s.constants.CMD_EDIT_PRICE_SAVE)
		return nil
	}

	wish[0].Price = message.Text
	if err := s.repositories.Wish.Save(ctx, wish[0]); err != nil {
		return err
	}

	s.clients.Telegram.Reply(ctx, s.constants.PRICE_SET, update)
	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), "")

	return nil
}

func (s *Service) EditLinkHandler(ctx context.Context, update *telegram.Update) error {
	s.clients.Telegram.Reply(ctx, s.constants.ENTER_LINK, update)

	s.clients.Cache.AppendText(ctx, update.GetChatIdStr(), update.CallbackQuery.Data)

	refreshMessageId := update.CallbackQuery.Message.GetMessageIdStr()
	s.clients.Cache.AppendText(ctx, update.GetChatIdStr(), refreshMessageId)

	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), s.constants.CMD_EDIT_LINK_SAVE)

	return nil
}

func (s *Service) SaveEditLinkHandler(ctx context.Context, update *telegram.Update) error {
	message := update.GetMessage()
	texts, err := s.clients.Cache.GetTexts(ctx, update.GetChatIdStr())
	if err != nil {
		return err
	}
	params := s.builders.CallbackDataBuilder.FromString(texts[0])
	wishId := params.ID

	link := message.Text
	if len(link) > s.constants.WISH_LINK_MAX_LEN {
		s.clients.Telegram.Reply(ctx, fmt.Sprintf(s.constants.LINK_TOO_LONG_TEMPLATE, s.constants.WISH_LINK_MAX_LEN), update)
		s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), s.constants.CMD_EDIT_LINK_SAVE)
		return nil
	}

	parsedURL, err := url.Parse(link)
	if err != nil || parsedURL.Host == "" {
		s.clients.Telegram.Reply(ctx, s.constants.LINK_INVALID_FORMAT, update)
		s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), s.constants.CMD_EDIT_LINK_SAVE)
		return nil
	}

	info, err := s.getCertificateInfo(link)
	if err != nil {
		s.clients.Telegram.Reply(ctx, s.constants.LINK_UNTRUSTED_SITE, update)
		s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), s.constants.CMD_EDIT_LINK_SAVE)
		return nil
	}

	s.logger.Debug("certificate check", "info", info)

	wish, err := s.repositories.Wish.Filter(ctx, &entities.Wish{BaseFields: entities.BaseFields{ID: wishId}})
	if len(wish) == 0 || err != nil {
		s.clients.Telegram.Reply(ctx, s.constants.WISH_NOT_FOUND, update)
		return err
	}
	wish[0].Link = link
	if err := s.repositories.Wish.Save(ctx, wish[0]); err != nil {
		return err
	}

	chatId := update.GetMessage().GetChatIdStr()

	// delete message with link - it's too large
	s.clients.Telegram.DeleteMessage(ctx, chatId, message.GetMessageIdStr())

	s.clients.Telegram.Reply(ctx, s.constants.LINK_SET, update)
	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), "")

	return nil
}

func (s *Service) EditWishNameHandler(ctx context.Context, update *telegram.Update) error {
	s.clients.Telegram.Reply(ctx, s.constants.ENTER_NEW_WISH_NAME, update)

	s.clients.Cache.AppendText(ctx, update.GetChatIdStr(), update.CallbackQuery.Data)

	refreshMessageId := update.CallbackQuery.Message.GetMessageIdStr()
	s.clients.Cache.AppendText(ctx, update.GetChatIdStr(), refreshMessageId)

	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), s.constants.CMD_EDIT_WISH_NAME_SAVE)

	return nil
}

func (s *Service) SaveEditWishNameHandler(ctx context.Context, update *telegram.Update) error {
	message := update.GetMessage()
	texts, err := s.clients.Cache.GetTexts(ctx, update.GetChatIdStr())
	if err != nil {
		return err
	}
	params := s.builders.CallbackDataBuilder.FromString(texts[0])
	wishId := params.ID

	_, err = s.validateWishName(message.Text)
	if err != nil {
		s.clients.Telegram.Reply(ctx, s.constants.WISH_NAME_INVALID_CHARS, update)
		s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), s.constants.CMD_EDIT_WISH_NAME_SAVE)
		return nil
	}

	wish, err := s.repositories.Wish.Filter(ctx, &entities.Wish{BaseFields: entities.BaseFields{ID: wishId}})
	if err != nil {
		s.logger.Error(
			"SaveEditWishNameHandler",
			"details", err.Error(),
			"chatId", update.GetMessage().GetChatIdStr(),
		)
		s.clients.Telegram.Reply(ctx, s.constants.WISH_NOT_FOUND, update)
		return err
	}

	wish[0].Name = message.Text
	err = s.repositories.Wish.Save(ctx, wish[0])
	if err != nil {
		return err
	}

	s.clients.Telegram.Reply(ctx, s.constants.WISH_NAME_CHANGED, update)
	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), "")

	return nil
}

func (s *Service) getCertificateInfo(urlStr string) (string, error) {
	if !strings.HasPrefix(urlStr, s.constants.HTTPS_PREFIX) {
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
