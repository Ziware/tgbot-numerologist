package communicate

import (
	"tgbot-numerologist/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LogMessage(bot *tgbotapi.BotAPI, req *tgbotapi.Message) {
	utils.Log("[%s] %s", req.From.UserName, req.Text)
}

func SendCommon(bot *tgbotapi.BotAPI, req *tgbotapi.Message) bool {
	msgConf := tgbotapi.NewMessage(req.Chat.ID, utils.HelpMessage)
	msg, err := bot.Send(msgConf)
	if err != nil {
		utils.Log("Error sending message: %v", err)
		return false
	} else {
		LogMessage(bot, &msg)
		return true
	}
}

func SendText(bot *tgbotapi.BotAPI, chatId int64, text string) bool {
	msgConf := tgbotapi.NewMessage(chatId, text)
	msg, err := bot.Send(msgConf)
	if err != nil {
		utils.Log("Error sending message: %v", err)
		return false
	} else {
		LogMessage(bot, &msg)
		return true
	}
}

func SendMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) bool {
	m, err := bot.Send(msg)
	if err != nil {
		utils.Log("Error sending message: %v", err)
		return false
	} else {
		LogMessage(bot, &m)
		return true
	}
}

func SendError(bot *tgbotapi.BotAPI, chatId int64, err error) bool {
	msgConf := tgbotapi.NewMessage(chatId, err.Error())
	msg, err := bot.Send(msgConf)
	if err != nil {
		utils.Log("Error sending message: %v", err)
		return false
	} else {
		LogMessage(bot, &msg)
		return true
	}
}
