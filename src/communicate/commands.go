package communicate

import (
	"errors"

	"tgbot-numerologist/database"
	"tgbot-numerologist/objects"
	"tgbot-numerologist/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	SendText(bot, message.Chat.ID, utils.HelpMessage)
}

func HandleProfile(bot *tgbotapi.BotAPI, message *tgbotapi.Message, profile *objects.Profile) {
	msgText := objects.FormatProfileMessage(profile)
	msg := tgbotapi.NewMessage(profile.ChatID, msgText)
	msg.ReplyMarkup = objects.ProfileKeyboard()
	msg.ParseMode = "Markdown"

	SendMessage(bot, &msg)
}

func HandleEditMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, profile *objects.Profile) {
	chatID := message.Chat.ID

	switch profile.EditingField {
	case "edit_name":
		profile.Name = message.Text
	case "edit_surname":
		profile.Surname = message.Text
	case "edit_birthdate":
		birthDate, err := objects.ParseDate(message.Text)
		if err == nil {
			profile.BirthDate = birthDate
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, "Неправильный формат даты, попробуйте ещё раз"))
			return
		}
	case "edit_bio":
		profile.Bio = message.Text
	case "edit_workplace":
		profile.WorkPlace = message.Text
	case "edit_studyplace":
		profile.StudyPlace = message.Text
	case "edit_hobby":
		profile.Hobby = message.Text
	default:
		bot.Send(tgbotapi.NewMessage(chatID, "Чтобы отредактировать поле, используйте кнопки в меню профиля"))
		return
	}

	profile.EditingField = ""

	msgText := "Ваш профиль обновлен:\n" + objects.FormatProfileMessage(profile)
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = objects.ProfileKeyboard()
	msg.ParseMode = "Markdown"

	err := database.SaveProfileToRedis(profile)
	if err != nil {
		utils.Log("error on save profile when edit: %s", err.Error())
		SendError(bot, message.Chat.ID, errors.New(utils.ErrGotSomeProblems))
		return
	}

	SendMessage(bot, &msg)
}

func HandlePredictions(bot *tgbotapi.BotAPI, message *tgbotapi.Message, profile *objects.Profile) {
	SendText(bot, message.Chat.ID, "Предсказаний пока что нет (")
}

func HandleStop(bot *tgbotapi.BotAPI, message *tgbotapi.Message, profile *objects.Profile) {
	profile.EditingField = ""
	err := database.SaveProfileToRedis(profile)
	if err != nil {
		utils.Log("error on save profile when edit: %s", err.Error())
		SendError(bot, message.Chat.ID, errors.New(utils.ErrGotSomeProblems))
		return
	}
	msgText := "Ввод отменен:\n" + objects.FormatProfileMessage(profile)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = objects.ProfileKeyboard()
	msg.ParseMode = "Markdown"
	SendMessage(bot, &msg)
}
