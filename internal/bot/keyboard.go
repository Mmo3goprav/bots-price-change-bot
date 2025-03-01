package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func newTextMsg(chatID int64, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = parseModHTML

	return msg
}

func newChartKeyboard(chatId int64, text string) tgbotapi.MessageConfig {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Add Chart"),
			tgbotapi.NewKeyboardButton("Remove Chart"),
		),
	)

	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = keyboard

	msg.ParseMode = parseModHTML

	return msg
}
