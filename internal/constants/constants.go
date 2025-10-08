package constants

import "github.com/meehighlov/grats/internal/config"

type Constants struct {
	ERROR_MESSAGE string

	// User messages
	GREETING_FRIEND   string
	GREETING_TEMPLATE string
	HELLO_AGAIN       string

	// Wish management
	WISH_LIMIT_REACHED_TEMPLATE string
	ENTER_WISH_NAME             string
	WISH_LIMIT_INFO_TEMPLATE    string
	WISH_NAME_TOO_LONG_TEMPLATE string
	WISH_NAME_INVALID_CHARS     string
	WISH_ADDED                  string
	WISH_DELETED                string
	WISH_NOT_FOUND              string
	WISH_WAS_DELETED            string
	WISH_ALREADY_BOOKED         string
	WISH_REMOVED_TRY_REFRESH    string

	// Price management
	ENTER_PRICE          string
	PRICE_INVALID_FORMAT string
	PRICE_SET            string

	// Link management
	ENTER_LINK             string
	LINK_TOO_LONG_TEMPLATE string
	LINK_INVALID_FORMAT    string
	LINK_UNTRUSTED_SITE    string
	LINK_SET               string

	// Name editing
	ENTER_NEW_WISH_NAME string
	WISH_NAME_CHANGED   string

	// Wishlist
	MY_WISHLIST              string
	WISHLIST_EMPTY           string
	WISHLIST_HEADER_TEMPLATE string
	SHARE_WISHLIST_MESSAGE   string
	FAILED_TO_LOAD_WISHES    string
	DEFAULT_WISHLIST_NAME    string
	WISHLIST_CREATION_ERROR  string

	// Buttons
	BTN_SHARE            string
	BTN_COPY_LINK        string
	BTN_BACK_TO_WISHLIST string
	BTN_REFRESH          string
	BTN_CANCEL_BOOKING   string
	BTN_BOOK_WISH        string
	BTN_BACK             string
	BTN_DELETE           string
	BTN_EDIT_NAME        string
	BTN_EDIT_LINK        string
	BTN_EDIT_PRICE       string
	BTN_OPEN_WISH        string
	BTN_NEW_WISH         string
	BTN_WISH_LIST        string
	BTN_ADD_WISH         string
	BTN_PREVIOUS         string
	BTN_SHARE_LIST       string
	BTN_WRITE            string
	BTN_CANCEL           string

	// Status messages
	STATUS_WISH_AVAILABLE       string
	STATUS_WISH_BOOKED_BY_YOU   string
	STATUS_WISH_BOOKED_BY_OTHER string

	// Other
	VIEW_ON_SITE  string
	CURRENCY_RUB  string
	LOCKED_PREFIX string
	HTTPS_PREFIX  string

	// Delete confirmation
	DELETE_WISH_CONFIRMATION_TEMPLATE string

	// Wishlist sharing
	MY_WISHLIST_SHARE_TITLE      string
	SHARE_WISHLIST_LINK_TEMPLATE string
	SHARED_LIST_ID_PREFIX        string

	// Support
	SUPPORT_REQUEST_MESSAGE  string
	SUPPORT_SEND_MESSAGE     string
	SUPPORT_MESSAGE_SENT     string
	SUPPORT_MESSAGE_TOO_LONG string
	SUPPORT_MESSAGE_TEMPLATE string
	SUPPORT_REPLY_TEMPLATE   string
	SUPPORT_CHAT_ID_PREFIX   string

	// Commands for callback data
	CMD_LIST                string
	CMD_WISH_INFO           string
	CMD_ADD_TO_WISH         string
	CMD_DELETE_WISH         string
	CMD_CONFIRM_DELETE_WISH string
	CMD_EDIT_WISH_NAME      string
	CMD_EDIT_LINK           string
	CMD_EDIT_PRICE          string
	CMD_TOGGLE_WISH_LOCK    string
	CMD_SHARE_WISH_LIST     string
	CMD_SHOW_SWL            string
	CMD_SHOW_SWI            string
	CMD_WISH_LIST           string
	CMD_SUPPORT_WRITE       string
	CMD_SUPPORT_CANCEL      string

	// Handler commands
	CMD_ADD_SAVE_WISH       string
	CMD_EDIT_PRICE_SAVE     string
	CMD_EDIT_LINK_SAVE      string
	CMD_EDIT_WISH_NAME_SAVE string
	CMD_SUPPORT_SEND        string

	// Bot commands
	CMD_START    string
	CMD_WISHLIST string
	CMD_SUPPORT  string

	// Numeric constants
	WISH_LIMIT_FOR_USER int
	WISH_NAME_MAX_LEN   int
	LIST_START_OFFSET   int
	WISH_LINK_MAX_LEN   int
	LIST_DEFAULT_OFFSET string
}

func New(cfg *config.Config) *Constants {
	return &Constants{
		ERROR_MESSAGE: "Возникла непредвиденная ошибка😔",

		// User messages
		GREETING_FRIEND:   "друг",
		GREETING_TEMPLATE: "Привет, %s👋",
		HELLO_AGAIN:       "Снова привет👋",

		// Wish management
		WISH_LIMIT_REACHED_TEMPLATE: "Достигнут лимит желаний👉👈 Максимальное количество желаний для одного пользователя: %d",
		ENTER_WISH_NAME:             "✨Введите название желания\n",
		WISH_NAME_TOO_LONG_TEMPLATE: "Слишком большое имя, максимум - %d символов, попробуйте снова",
		WISH_NAME_INVALID_CHARS:     "Имя желания содержит недопустимые символы, попробуйте использовать только цифры и буквы",
		WISH_ADDED:                  "Желание добавлено 💾",
		WISH_DELETED:                "Желание удалено 🗑️",
		WISH_NOT_FOUND:              "Не смог найти желание",
		WISH_WAS_DELETED:            "Желание было удалено, попробуйте обновить список",
		WISH_ALREADY_BOOKED:         "Кто-то уже забронировал это желание, попробуйте обновить список",
		WISH_REMOVED_TRY_REFRESH:    "Видимо, желание было удалено🤔 Попробуйте обновить список",

		// Price management
		ENTER_PRICE:          "Введите цену в рублях",
		PRICE_INVALID_FORMAT: "Не могу распознать цену, пример 👉 1000, 400.5",
		PRICE_SET:            "Цена установлена 💾",

		// Link management
		ENTER_LINK:             "Введите ссылку",
		LINK_TOO_LONG_TEMPLATE: "Ссылка не может быть длиннее %d символов",
		LINK_INVALID_FORMAT:    "Не могу распознать ссылку, введите ссылку в формате https://example.com",
		LINK_UNTRUSTED_SITE:    "Я не доверяю этому сайту😔 попробуй другую ссылку",
		LINK_SET:               "Ссылка установлена 💾",

		// Name editing
		ENTER_NEW_WISH_NAME: "Введите новое название желания",
		WISH_NAME_CHANGED:   "Название желания изменено 💾",

		// Wishlist
		MY_WISHLIST:              "✨Мой вишлист",
		WISHLIST_EMPTY:           "✨Список пуст",
		WISHLIST_HEADER_TEMPLATE: "✨Вишлист %s",
		SHARE_WISHLIST_MESSAGE:   "Поделитесь своим вишлистом!\n\n- При переходе по ссылке откроется чат с grats и вишлист будет прислан в виде нового сообщения\n- Если пользователь ранее не использовал grats, то ему потребуется лишь нажать start",
		FAILED_TO_LOAD_WISHES:    "Не удалось загрузить желания пользователя😔",
		DEFAULT_WISHLIST_NAME:    "Мой wishlist",
		WISHLIST_CREATION_ERROR:  "Возникла непредвиденная ошибка при создании первого списка желаний, над этим уже работают😔",

		// Buttons
		BTN_SHARE:            "📤 поделиться",
		BTN_COPY_LINK:        "🔗 ссылка",
		BTN_BACK_TO_WISHLIST: "⬅️ к списку желаний",
		BTN_REFRESH:          "🔄",
		BTN_CANCEL_BOOKING:   "✖️ отменить бронь",
		BTN_BOOK_WISH:        "🎁 забронировать",
		BTN_BACK:             "⬅️ назад",
		BTN_DELETE:           "🗑️ удалить",
		BTN_EDIT_NAME:        "✏️ название",
		BTN_EDIT_LINK:        "✏️ ссылка",
		BTN_EDIT_PRICE:       "✏️ цена",
		BTN_OPEN_WISH:        "📂 открыть желание",
		BTN_NEW_WISH:         "➕ новое желание",
		BTN_WISH_LIST:        "📋 список желаний",
		BTN_ADD_WISH:         "➕",
		BTN_SHARE_LIST:       "🛜",
		BTN_PREVIOUS:         "⬅️",
		BTN_WRITE:            "✍️ написать",
		BTN_CANCEL:           "❌ отмена",

		// Status messages
		STATUS_WISH_AVAILABLE:       "🟢 желание пока не выбрали",
		STATUS_WISH_BOOKED_BY_YOU:   "🎁 вы забронировали это желание",
		STATUS_WISH_BOOKED_BY_OTHER: "🎁 кто-то забронировал это желание",

		// Other
		VIEW_ON_SITE:  "смотреть на сайте",
		CURRENCY_RUB:  "(RUB)",
		LOCKED_PREFIX: "🔒 ",
		HTTPS_PREFIX:  "https://",

		// Delete confirmation
		DELETE_WISH_CONFIRMATION_TEMPLATE: "Удалить желание %s?",

		// Wishlist sharing
		MY_WISHLIST_SHARE_TITLE:      "Мой wishlist✨",
		SHARE_WISHLIST_LINK_TEMPLATE: "https://t.me/%s?start=wl%s",
		SHARED_LIST_ID_PREFIX:        "wl",

		// Support
		SUPPORT_REQUEST_MESSAGE:  "💬 Оставьте отзыв или задайте вопрос",
		SUPPORT_SEND_MESSAGE:     "📝 Отправьте сообщение, оно будет рассмотрено поддержкой в ближайшее время, ответ придет в диалог с ботом",
		SUPPORT_MESSAGE_SENT:     "✅ Ваше сообщение отправлено! Ответ придет в этот чат",
		SUPPORT_MESSAGE_TOO_LONG: "Сообщение слишком длинное, максимум 2000 символов",
		SUPPORT_MESSAGE_TEMPLATE: "chatid:%s\n📨 Новое сообщение от пользователя %s (ID: %s)\n\n%s",
		SUPPORT_REPLY_TEMPLATE:   "💬 Ответ от поддержки:\n\n%s",
		SUPPORT_CHAT_ID_PREFIX:   "chatid:",

		// Commands for callback data
		CMD_LIST:                "list",
		CMD_WISH_INFO:           "wish_info",
		CMD_ADD_TO_WISH:         "add_to_wish",
		CMD_DELETE_WISH:         "delete_wish",
		CMD_CONFIRM_DELETE_WISH: "confirm_delete_wish",
		CMD_EDIT_WISH_NAME:      "edit_wish_name",
		CMD_EDIT_LINK:           "edit_link",
		CMD_EDIT_PRICE:          "edit_price",
		CMD_TOGGLE_WISH_LOCK:    "toggle_wish_lock",
		CMD_SHARE_WISH_LIST:     "share_wish_list",
		CMD_SHOW_SWL:            "show_swl",
		CMD_SHOW_SWI:            "show_swi",
		CMD_WISH_LIST:           "wish_list",
		CMD_SUPPORT_WRITE:       "support_write",
		CMD_SUPPORT_CANCEL:      "support_cancel",

		// Handler commands
		CMD_ADD_SAVE_WISH:       "add_save_wish",
		CMD_EDIT_PRICE_SAVE:     "edit_price_save",
		CMD_EDIT_LINK_SAVE:      "edit_link_save",
		CMD_EDIT_WISH_NAME_SAVE: "edit_wish_name_save",
		CMD_SUPPORT_SEND:        "support_send",

		// Bot commands
		CMD_START:    "/start",
		CMD_WISHLIST: "/wishlist",
		CMD_SUPPORT:  "/support",

		// Numeric constants
		WISH_LIMIT_FOR_USER: 50,
		WISH_NAME_MAX_LEN:   100,
		LIST_START_OFFSET:   0,
		WISH_LINK_MAX_LEN:   500,
		LIST_DEFAULT_OFFSET: "0",
	}
}
