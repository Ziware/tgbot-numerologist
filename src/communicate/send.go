package communicate

import (
	"tgbot-numerologist/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LogMessage(bot *tgbotapi.BotAPI, req *tgbotapi.Message) {
	utils.Log("[%s] %s", req.From.UserName, req.Text)
}

func SendCommon(bot *tgbotapi.BotAPI, req *tgbotapi.Message) {
	msgConf := tgbotapi.NewMessage(req.Chat.ID, utils.HelpMessage)
	msg, err := bot.Send(msgConf)
	if err != nil {
		utils.Log("Error sending message: %v", err)
	} else {
		LogMessage(bot, &msg)
	}
}

func SendText(bot *tgbotapi.BotAPI, chatId int64, text string) {
	msgConf := tgbotapi.NewMessage(chatId, text)
	msg, err := bot.Send(msgConf)
	if err != nil {
		utils.Log("Error sending message: %v", err)
	} else {
		LogMessage(bot, &msg)
	}
}

func SendMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) {
	m, err := bot.Send(msg)
	if err != nil {
		utils.Log("Error sending message: %v", err)
	} else {
		LogMessage(bot, &m)
	}
}

func SendError(bot *tgbotapi.BotAPI, chatId int64, err error) {
	msgConf := tgbotapi.NewMessage(chatId, err.Error())
	msg, err := bot.Send(msgConf)
	if err != nil {
		utils.Log("Error sending message: %v", err)
	} else {
		LogMessage(bot, &msg)
	}
}
