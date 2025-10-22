package config

type Constants struct {
	ERROR_MESSAGE string `env:"ERROR_MESSAGE" env-default:"Возникла непредвиденная ошибка😔"`

	// User messages
	GREETING_FRIEND   string `env:"GREETING_FRIEND" env-default:"друг"`
	GREETING_TEMPLATE string `env:"GREETING_TEMPLATE" env-default:"Привет, %s👋"`
	HELLO_AGAIN       string `env:"HELLO_AGAIN" env-default:"Снова привет👋"`

	// Wish management
	WISH_LIMIT_REACHED_TEMPLATE string `env:"WISH_LIMIT_REACHED_TEMPLATE" env-default:"Достигнут лимит желаний👉👈 Максимальное количество желаний для одного пользователя: %d"`
	ENTER_WISH_NAME             string `env:"ENTER_WISH_NAME" env-default:"✨Введите название желания\n"`
	WISH_LIMIT_INFO_TEMPLATE    string `env:"WISH_LIMIT_INFO_TEMPLATE" env-default:""`
	WISH_NAME_TOO_LONG_TEMPLATE string `env:"WISH_NAME_TOO_LONG_TEMPLATE" env-default:"Слишком большое имя, максимум - %d символов, попробуйте снова"`
	WISH_NAME_INVALID_CHARS     string `env:"WISH_NAME_INVALID_CHARS" env-default:"Имя желания содержит недопустимые символы, попробуйте использовать только цифры и буквы"`
	WISH_ADDED                  string `env:"WISH_ADDED" env-default:"Желание добавлено 💾"`
	WISH_DELETED                string `env:"WISH_DELETED" env-default:"Желание удалено 🗑️"`
	WISH_NOT_FOUND              string `env:"WISH_NOT_FOUND" env-default:"Не смог найти желание"`
	WISH_WAS_DELETED            string `env:"WISH_WAS_DELETED" env-default:"Желание было удалено, попробуйте обновить список"`
	WISH_ALREADY_BOOKED         string `env:"WISH_ALREADY_BOOKED" env-default:"Кто-то уже забронировал это желание, попробуйте обновить список"`
	WISH_REMOVED_TRY_REFRESH    string `env:"WISH_REMOVED_TRY_REFRESH" env-default:"Видимо, желание было удалено🤔 Попробуйте обновить список"`

	// Price management
	ENTER_PRICE          string `env:"ENTER_PRICE" env-default:"Введите цену в рублях"`
	PRICE_INVALID_FORMAT string `env:"PRICE_INVALID_FORMAT" env-default:"Не могу распознать цену, пример 👉 1000, 400.5"`
	PRICE_SET            string `env:"PRICE_SET" env-default:"Цена установлена 💾"`

	// Link management
	ENTER_LINK             string `env:"ENTER_LINK" env-default:"Введите ссылку"`
	LINK_TOO_LONG_TEMPLATE string `env:"LINK_TOO_LONG_TEMPLATE" env-default:"Ссылка не может быть длиннее %d символов"`
	LINK_INVALID_FORMAT    string `env:"LINK_INVALID_FORMAT" env-default:"Не могу распознать ссылку, введите ссылку в формате https://example.com"`
	LINK_UNTRUSTED_SITE    string `env:"LINK_UNTRUSTED_SITE" env-default:"Я не доверяю этому сайту😔 попробуй другую ссылку"`
	LINK_SET               string `env:"LINK_SET" env-default:"Ссылка установлена 💾"`
	LINK_DELETED           string `env:"LINK_DELETED" env-default:"ссылка удалена"`

	// Name editing
	ENTER_NEW_WISH_NAME string `env:"ENTER_NEW_WISH_NAME" env-default:"Введите новое название желания"`
	WISH_NAME_CHANGED   string `env:"WISH_NAME_CHANGED" env-default:"Название желания изменено 💾"`

	// Wishlist
	MY_WISHLIST              string `env:"MY_WISHLIST" env-default:"✨Мой вишлист"`
	WISHLIST_EMPTY           string `env:"WISHLIST_EMPTY" env-default:"✨Список пуст"`
	WISHLIST_HEADER_TEMPLATE string `env:"WISHLIST_HEADER_TEMPLATE" env-default:"✨Вишлист %s"`
	SHARE_WISHLIST_MESSAGE   string `env:"SHARE_WISHLIST_MESSAGE" env-default:"Поделитесь своим вишлистом!\n\n- При переходе по ссылке откроется чат с grats и вишлист будет прислан в виде нового сообщения\n- Если пользователь ранее не использовал grats, то ему потребуется лишь нажать start"`
	FAILED_TO_LOAD_WISHES    string `env:"FAILED_TO_LOAD_WISHES" env-default:"Не удалось загрузить желания пользователя😔"`
	DEFAULT_WISHLIST_NAME    string `env:"DEFAULT_WISHLIST_NAME" env-default:"Мой wishlist"`
	WISHLIST_CREATION_ERROR  string `env:"WISHLIST_CREATION_ERROR" env-default:"Возникла непредвиденная ошибка при создании первого списка желаний, над этим уже работают😔"`

	// Buttons
	BTN_SHARE            string `env:"BTN_SHARE" env-default:"📤 поделиться"`
	BTN_COPY_LINK        string `env:"BTN_COPY_LINK" env-default:"🔗 ссылка"`
	BTN_BACK_TO_WISHLIST string `env:"BTN_BACK_TO_WISHLIST" env-default:"⬅️ к списку желаний"`
	BTN_REFRESH          string `env:"BTN_REFRESH" env-default:"🔄"`
	BTN_CANCEL_BOOKING   string `env:"BTN_CANCEL_BOOKING" env-default:"✖️ отменить бронь"`
	BTN_BOOK_WISH        string `env:"BTN_BOOK_WISH" env-default:"🎁 забронировать"`
	BTN_BACK             string `env:"BTN_BACK" env-default:"⬅️ назад"`
	BTN_DELETE           string `env:"BTN_DELETE" env-default:"🗑️ удалить"`
	BTN_EDIT_NAME        string `env:"BTN_EDIT_NAME" env-default:"✏️ название"`
	BTN_EDIT_LINK        string `env:"BTN_EDIT_LINK" env-default:"✏️ ссылка"`
	BTN_EDIT_PRICE       string `env:"BTN_EDIT_PRICE" env-default:"✏️ цена"`
	BTN_OPEN_WISH        string `env:"BTN_OPEN_WISH" env-default:"📂 открыть желание"`
	BTN_NEW_WISH         string `env:"BTN_NEW_WISH" env-default:"➕ новое желание"`
	BTN_WISH_LIST        string `env:"BTN_WISH_LIST" env-default:"📋 список желаний"`
	BTN_ADD_WISH         string `env:"BTN_ADD_WISH" env-default:"➕"`
	BTN_PREVIOUS         string `env:"BTN_PREVIOUS" env-default:"⬅️"`
	BTN_SHARE_LIST       string `env:"BTN_SHARE_LIST" env-default:"🛜"`
	BTN_WRITE            string `env:"BTN_WRITE" env-default:"✍️ написать"`
	BTN_CANCEL           string `env:"BTN_CANCEL" env-default:"❌ отмена"`
	BTN_DELETE_LINK      string `env:"BTN_DELETE_LINK" env-default:"удалить ссылку 🗑️"`

	// Status messages
	STATUS_WISH_AVAILABLE       string `env:"STATUS_WISH_AVAILABLE" env-default:"🟢 желание пока не выбрали"`
	STATUS_WISH_BOOKED_BY_YOU   string `env:"STATUS_WISH_BOOKED_BY_YOU" env-default:"🎁 вы забронировали это желание"`
	STATUS_WISH_BOOKED_BY_OTHER string `env:"STATUS_WISH_BOOKED_BY_OTHER" env-default:"🎁 кто-то забронировал это желание"`

	// Other
	VIEW_ON_SITE  string `env:"VIEW_ON_SITE" env-default:"смотреть на сайте"`
	CURRENCY_RUB  string `env:"CURRENCY_RUB" env-default:"(RUB)"`
	LOCKED_PREFIX string `env:"LOCKED_PREFIX" env-default:"🔒 "`
	HTTPS_PREFIX  string `env:"HTTPS_PREFIX" env-default:"https://"`

	// Delete confirmation
	DELETE_WISH_CONFIRMATION_TEMPLATE string `env:"DELETE_WISH_CONFIRMATION_TEMPLATE" env-default:"Удалить желание %s?"`

	// Wishlist sharing
	MY_WISHLIST_SHARE_TITLE      string `env:"MY_WISHLIST_SHARE_TITLE" env-default:"Мой wishlist✨"`
	SHARE_WISHLIST_LINK_TEMPLATE string `env:"SHARE_WISHLIST_LINK_TEMPLATE" env-default:"https://t.me/%s?start=wl%s"`
	SHARED_LIST_ID_PREFIX        string `env:"SHARED_LIST_ID_PREFIX" env-default:"wl"`

	// Support
	SUPPORT_REQUEST_MESSAGE  string `env:"SUPPORT_REQUEST_MESSAGE" env-default:"💬 Оставьте отзыв или задайте вопрос"`
	SUPPORT_SEND_MESSAGE     string `env:"SUPPORT_SEND_MESSAGE" env-default:"📝 Отправьте сообщение, оно будет рассмотрено поддержкой в ближайшее время, ответ придет в диалог с ботом"`
	SUPPORT_MESSAGE_SENT     string `env:"SUPPORT_MESSAGE_SENT" env-default:"✅ Ваше сообщение отправлено! Ответ придет в этот чат"`
	SUPPORT_MESSAGE_TOO_LONG string `env:"SUPPORT_MESSAGE_TOO_LONG" env-default:"Сообщение слишком длинное, максимум 2000 символов"`
	SUPPORT_MESSAGE_TEMPLATE string `env:"SUPPORT_MESSAGE_TEMPLATE" env-default:"chatid:%s\n📨 Новое сообщение от пользователя %s (ID: %s)\n\n%s"`
	SUPPORT_REPLY_TEMPLATE   string `env:"SUPPORT_REPLY_TEMPLATE" env-default:"💬 Ответ от поддержки:\n\n%s"`
	SUPPORT_CHAT_ID_PREFIX   string `env:"SUPPORT_CHAT_ID_PREFIX" env-default:"chatid:"`

	// Commands for callback data
	CMD_LIST                string `env:"CMD_LIST" env-default:"list"`
	CMD_WISH_INFO           string `env:"CMD_WISH_INFO" env-default:"wish_info"`
	CMD_ADD_TO_WISH         string `env:"CMD_ADD_TO_WISH" env-default:"add_to_wish"`
	CMD_DELETE_WISH         string `env:"CMD_DELETE_WISH" env-default:"d_wish"`
	CMD_CONFIRM_DELETE_WISH string `env:"CMD_CONFIRM_DELETE_WISH" env-default:"c_delete_wish"`
	CMD_EDIT_WISH_NAME      string `env:"CMD_EDIT_WISH_NAME" env-default:"ewn"`
	CMD_EDIT_LINK           string `env:"CMD_EDIT_LINK" env-default:"ewl"`
	CMD_EDIT_PRICE          string `env:"CMD_EDIT_PRICE" env-default:"ewp"`
	CMD_TOGGLE_WISH_LOCK    string `env:"CMD_TOGGLE_WISH_LOCK" env-default:"toggle_wish_lock"`
	CMD_SHARE_WISH_LIST     string `env:"CMD_SHARE_WISH_LIST" env-default:"share_wl"`
	CMD_SHOW_SWL            string `env:"CMD_SHOW_SWL" env-default:"show_swl"`
	CMD_SHOW_SWI            string `env:"CMD_SHOW_SWI" env-default:"show_swi"`
	CMD_SUPPORT_WRITE       string `env:"CMD_SUPPORT_WRITE" env-default:"support_write"`
	CMD_SUPPORT_CANCEL      string `env:"CMD_SUPPORT_CANCEL" env-default:"support_cancel"`
	CMD_DELETE_LINK         string `env:"CMD_DELETE_LINK" env-default:"delete_link"`

	// Handler commands
	CMD_ADD_SAVE_WISH       string `env:"CMD_ADD_SAVE_WISH" env-default:"add_save_wish"`
	CMD_EDIT_PRICE_SAVE     string `env:"CMD_EDIT_PRICE_SAVE" env-default:"ewp_save"`
	CMD_EDIT_LINK_SAVE      string `env:"CMD_EDIT_LINK_SAVE" env-default:"ewl_save"`
	CMD_EDIT_WISH_NAME_SAVE string `env:"CMD_EDIT_WISH_NAME_SAVE" env-default:"ewn_save"`
	CMD_SUPPORT_SEND        string `env:"CMD_SUPPORT_SEND" env-default:"support_send"`

	// Bot commands
	CMD_START    string `env:"CMD_START" env-default:"/start"`
	CMD_WISHLIST string `env:"CMD_WISHLIST" env-default:"/wishlist"`
	CMD_SUPPORT  string `env:"CMD_SUPPORT" env-default:"/support"`
	CMD_CANCEL   string `env:"CMD_CANCEL" env-default:"/cancel"`

	// Numeric constants
	WISH_LIMIT_FOR_USER int    `env:"WISH_LIMIT_FOR_USER" env-default:"50"`
	WISH_NAME_MAX_LEN   int    `env:"WISH_NAME_MAX_LEN" env-default:"100"`
	LIST_START_OFFSET   int    `env:"LIST_START_OFFSET" env-default:"0"`
	WISH_LINK_MAX_LEN   int    `env:"WISH_LINK_MAX_LEN" env-default:"500"`
	LIST_DEFAULT_OFFSET string `env:"LIST_DEFAULT_OFFSET" env-default:"0"`
}
