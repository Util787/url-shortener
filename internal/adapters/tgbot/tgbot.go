package tgbot

import (
	"fmt"
	"html"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ShortenerUsecase describes the subset of usecase methods we need in the bot.
// Accepts the concrete usecase produced by shortener.NewShortenerUsecase from main.
type ShortenerUsecase interface {
	SaveURL(longURL string) (string, error)
	GetRandomURL() (string, error)
	DeleteURL(id *string, longURL *string, shortURL *string) error
}

type Bot struct {
	botAPI  *tgbotapi.BotAPI
	usecase ShortenerUsecase

	// simple per-chat state to wait for user input after a button press
	awaitingMu sync.Mutex
	awaiting   map[int64]string // chatID -> action: "save" | "delete"
}

func NewBot(token string, uc ShortenerUsecase) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = false
	log.Printf("Бот авторизован как %s", botAPI.Self.UserName)

	return &Bot{
		botAPI:   botAPI,
		usecase:  uc,
		awaiting: make(map[int64]string),
	}, nil
}

// Start runs the update loop (blocking). Call it in a goroutine.
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
	}
}

func (b *Bot) setAwaiting(chatID int64, action string) {
	b.awaitingMu.Lock()
	defer b.awaitingMu.Unlock()
	if action == "" {
		delete(b.awaiting, chatID)
	} else {
		b.awaiting[chatID] = action
	}
}

func (b *Bot) getAwaiting(chatID int64) string {
	b.awaitingMu.Lock()
	defer b.awaitingMu.Unlock()
	return b.awaiting[chatID]
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	// If user is responding to a previous action, handle it
	if action := b.getAwaiting(chatID); action != "" && !message.IsCommand() {
		switch action {
		case "save":
			longURL := message.Text
			shortURL, err := b.usecase.SaveURL(longURL)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "Ошибка при сохранении: "+err.Error())
				b.botAPI.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "URL успешно сохранён с короткой ссылкой: "+shortURL)
				b.botAPI.Send(msg)
			}
			b.setAwaiting(chatID, "")
			return
		case "delete":
			short := message.Text
			if err := b.usecase.DeleteURL(nil, nil, &short); err != nil {
				msg := tgbotapi.NewMessage(chatID, "Ошибка при удалении: "+err.Error())
				b.botAPI.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "URL успешно удалён")
				b.botAPI.Send(msg)
			}
			b.setAwaiting(chatID, "")
			return
		}
	}

	// Handle commands
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(chatID, "Привет! Выберите действие:")

			btnSave := tgbotapi.NewInlineKeyboardButtonData("Сохранить URL", "save")
			btnDelete := tgbotapi.NewInlineKeyboardButtonData("Удалить URL", "delete")
			btnRandom := tgbotapi.NewInlineKeyboardButtonData("Случайный URL", "random")

			row := tgbotapi.NewInlineKeyboardRow(btnSave, btnDelete, btnRandom)
			markup := tgbotapi.NewInlineKeyboardMarkup(row)
			msg.ReplyMarkup = markup
			b.botAPI.Send(msg)
		default:
			msg := tgbotapi.NewMessage(chatID, "Неизвестная команда")
			b.botAPI.Send(msg)
		}
	}
}

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) {
	// Acknowledge callback
	b.botAPI.Request(tgbotapi.NewCallback(callback.ID, ""))

	chatID := callback.Message.Chat.ID

	switch callback.Data {
	case "save":
		b.setAwaiting(chatID, "save")
		msg := tgbotapi.NewMessage(chatID, "Отправьте длинную ссылку, которую хотите сохранить:")
		b.botAPI.Send(msg)
	case "delete":
		b.setAwaiting(chatID, "delete")
		msg := tgbotapi.NewMessage(chatID, "Отправьте короткую ссылку (например: abc123) для удаления:")
		b.botAPI.Send(msg)
	case "random":
		shortURL, err := b.usecase.GetRandomURL()
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка при получении случайного URL: "+err.Error())
			b.botAPI.Send(msg)
			return
		}
		// Send clickable link using HTML parse mode. Escape user data to be safe.
		esc := html.EscapeString(shortURL)
		text := fmt.Sprintf(`Случайный URL: <a href="%s">%s</a> `, esc, esc)
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = tgbotapi.ModeHTML
		b.botAPI.Send(msg)
	default:
		msg := tgbotapi.NewMessage(chatID, "Неизвестное действие")
		b.botAPI.Send(msg)
	}
}
