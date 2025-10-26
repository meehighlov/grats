package wish

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/models"
)

func (s *Service) EditPrice(ctx context.Context, update *telegram.Update) error {
	s.clients.Telegram.Reply(ctx, s.cfg.Constants.ENTER_PRICE, update)

	s.repositories.Cache.AppendText(ctx, update.GetChatIdStr(), update.CallbackQuery.Data)

	refreshMessageId := update.CallbackQuery.Message.GetMessageIdStr()
	s.repositories.Cache.AppendText(ctx, update.GetChatIdStr(), refreshMessageId)

	return nil
}

func (s *Service) SaveEditPrice(ctx context.Context, update *telegram.Update) error {
	var (
		wish *models.Wish
	)
	message := update.GetMessage()
	texts, err := s.repositories.Cache.GetTexts(ctx, update.GetChatIdStr())
	if err != nil {
		return err
	}
	params := s.builders.CallbackDataBuilder.FromString(texts[0])
	wishId := params.ID

	err = s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wish, err = s.repositories.Wish.Get(ctx, wishId)
		if err != nil {
			s.clients.Telegram.Reply(ctx, s.cfg.Constants.WISH_NOT_FOUND, update)
			return err
		}

		price, err := strconv.ParseFloat(message.Text, 64)
		if err != nil || price <= 0 {
			s.clients.Telegram.Reply(ctx, s.cfg.Constants.PRICE_INVALID_FORMAT, update)
			return errors.New("price is invalid format")
		}

		wish.Price = message.Text
		return s.repositories.Wish.Save(ctx, wish)
	})
	if err != nil {
		return err
	}

	s.clients.Telegram.Reply(ctx, s.cfg.Constants.PRICE_SET, update)

	return nil
}

func (s *Service) EditLink(ctx context.Context, update *telegram.Update) error {
	var (
		wish *models.Wish
	)
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)
	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wish, err = s.repositories.Wish.Get(ctx, params.ID)
		return err
	})
	if err != nil {
		return err
	}
	if wish.Link != "" {
		keyboard := s.builders.KeyboardBuilder.NewKeyboard()
		btnCallback := s.builders.CallbackDataBuilder.Build(wish.ID, s.cfg.Constants.CMD_DELETE_LINK, params.Offset)
		keyboard.AppendAsLine(
			keyboard.NewButton(s.cfg.Constants.BTN_DELETE_LINK, btnCallback.String()),
		)
		s.clients.Telegram.Reply(
			ctx,
			s.cfg.Constants.ENTER_LINK,
			update,
			telegram.WithReplyMurkup(keyboard.Murkup()),
		)
	} else {
		s.clients.Telegram.Reply(ctx, s.cfg.Constants.ENTER_LINK, update)
	}

	s.repositories.Cache.AppendText(ctx, update.GetChatIdStr(), update.CallbackQuery.Data)

	refreshMessageId := update.CallbackQuery.Message.GetMessageIdStr()
	s.repositories.Cache.AppendText(ctx, update.GetChatIdStr(), refreshMessageId)

	return nil
}

func (s *Service) SaveEditLink(ctx context.Context, update *telegram.Update) error {
	var (
		wish *models.Wish
	)
	message := update.GetMessage()
	texts, err := s.repositories.Cache.GetTexts(ctx, update.GetChatIdStr())
	if err != nil {
		return err
	}
	params := s.builders.CallbackDataBuilder.FromString(texts[0])
	wishId := params.ID

	link := message.Text
	if len(link) > s.cfg.Constants.WISH_LINK_MAX_LEN {
		s.clients.Telegram.Reply(ctx, fmt.Sprintf(s.cfg.Constants.LINK_TOO_LONG_TEMPLATE, s.cfg.Constants.WISH_LINK_MAX_LEN), update)
		return errors.New("link is too long")
	}

	parsedURL, err := url.Parse(link)
	if err != nil || parsedURL.Host == "" {
		s.clients.Telegram.Reply(ctx, s.cfg.Constants.LINK_INVALID_FORMAT, update)
		return errors.New("link is invalid format")
	}

	info, err := s.getCertificateInfo(link)
	if err != nil {
		s.clients.Telegram.Reply(ctx, s.cfg.Constants.LINK_UNTRUSTED_SITE, update)
		return errors.New("link is untrusted site")
	}

	s.logger.Debug("certificate check", "info", info)

	err = s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wish, err = s.repositories.Wish.Get(ctx, wishId)
		if err != nil {
			s.clients.Telegram.Reply(ctx, s.cfg.Constants.WISH_NOT_FOUND, update)
			return err
		}
		wish.Link = link
		return s.repositories.Wish.Save(ctx, wish)
	})
	if err != nil {
		return err
	}

	chatId := update.GetMessage().GetChatIdStr()

	// delete message with link - it's too large
	s.clients.Telegram.DeleteMessage(ctx, chatId, message.GetMessageIdStr())

	s.clients.Telegram.Reply(ctx, s.cfg.Constants.LINK_SET, update)

	return nil
}

func (s *Service) DeleteLink(ctx context.Context, update *telegram.Update) error {
	var (
		wish *models.Wish
	)
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)
	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wish, err = s.repositories.Wish.Get(ctx, params.ID)
		if err != nil {
			return err
		}
		wish.Link = ""
		return s.repositories.Wish.Save(ctx, wish)
	})
	if err != nil {
		return err
	}

	s.clients.Telegram.Edit(ctx, s.cfg.Constants.LINK_DELETED, update)

	return nil
}

func (s *Service) EditWishName(ctx context.Context, update *telegram.Update) error {
	s.clients.Telegram.Reply(ctx, s.cfg.Constants.ENTER_NEW_WISH_NAME, update)

	s.repositories.Cache.AppendText(ctx, update.GetChatIdStr(), update.CallbackQuery.Data)

	refreshMessageId := update.CallbackQuery.Message.GetMessageIdStr()
	s.repositories.Cache.AppendText(ctx, update.GetChatIdStr(), refreshMessageId)

	return nil
}

func (s *Service) SaveEditWishName(ctx context.Context, update *telegram.Update) error {
	var (
		wish *models.Wish
	)
	message := update.GetMessage()
	texts, err := s.repositories.Cache.GetTexts(ctx, update.GetChatIdStr())
	if err != nil {
		return err
	}
	params := s.builders.CallbackDataBuilder.FromString(texts[0])
	wishId := params.ID

	_, err = s.validateWishName(message.Text)
	if err != nil {
		s.clients.Telegram.Reply(ctx, s.cfg.Constants.WISH_NAME_INVALID_CHARS, update)
		return errors.New("wish name contains invalid characters")
	}

	err = s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wish, err = s.repositories.Wish.Get(ctx, wishId)
		if err != nil {
			s.clients.Telegram.Reply(ctx, s.cfg.Constants.WISH_NOT_FOUND, update)
			return err
		}

		wish.Name = message.Text
		return s.repositories.Wish.Save(ctx, wish)
	})
	if err != nil {
		return err
	}

	s.clients.Telegram.Reply(ctx, s.cfg.Constants.WISH_NAME_CHANGED, update)

	return nil
}

func (s *Service) getCertificateInfo(urlStr string) (*x509.Certificate, error) {
	if !strings.HasPrefix(urlStr, s.cfg.Constants.HTTPS_PREFIX) {
		return nil, fmt.Errorf("URL must begun with https://")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	host := parsedURL.Host
	if !strings.Contains(host, ":") {
		host = host + ":443"
	}

	conn, err := tls.Dial("tcp", host, &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		return nil, fmt.Errorf("TLS connect error: %v", err)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certs found")
	}

	cert := certs[0]

	_ = fmt.Sprintf(
		"Субъект: %s\nИздатель: %s\nДействителен с: %s\nДействителен до: %s\nSAN: %v",
		cert.Subject,
		cert.Issuer,
		cert.NotBefore.Format("2006-01-02"),
		cert.NotAfter.Format("2006-01-02"),
		cert.DNSNames,
	)

	return cert, nil
}

func (s *Service) GetSiteName(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	host := parsedURL.Host
	if host == "" {
		return "", fmt.Errorf("cannot get hostname from link")
	}

	if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}

	host = strings.TrimPrefix(host, "www.")

	return host, nil
}
