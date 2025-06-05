package communicate

import (
	"errors"
	"fmt"

	"tgbot-numerologist/ai"
	"tgbot-numerologist/database"
	"tgbot-numerologist/objects"
	"tgbot-numerologist/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleIntro(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	SendText(bot, message.Chat.ID, utils.IntroMessage)
}

func HandleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	SendText(bot, message.Chat.ID, utils.HelpMessage)
}

func HandlePayment(bot *tgbotapi.BotAPI, message *tgbotapi.Message, profile *objects.Profile) {
	msgText := fmt.Sprintf(utils.PaymentMessage, profile.Quote, profile.Predictions)
	msg := tgbotapi.NewMessage(profile.ChatID, msgText)
	msg.ReplyMarkup = profile.GetPaymentKeyboard()
	msg.ParseMode = "Markdown"

	SendMessage(bot, &msg)
}

func HandlePayButton(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, profile *objects.Profile) {
	chatID := callbackQuery.Message.Chat.ID
	profile.Quote += 1
	err := database.SaveProfileToRedis(profile)
	if err != nil {
		utils.Log("error on save profile when edit: %s", err.Error())
		SendError(bot, callbackQuery.Message.Chat.ID, errors.New(utils.ErrGotSomeProblems))
		return
	}
	utils.Log("successfully saved data to redis")
	SendText(bot, chatID, "Квота увеличена")
	msgText := fmt.Sprintf(utils.PaymentMessage, profile.Quote, profile.Predictions)
	msg := tgbotapi.NewMessage(profile.ChatID, msgText)
	msg.ReplyMarkup = profile.GetPaymentKeyboard()
	msg.ParseMode = "Markdown"

	SendMessage(bot, &msg)
}

func HandleReset(bot *tgbotapi.BotAPI, message *tgbotapi.Message, profile *objects.Profile) {
	profile.ResetProfile()
	err := database.SaveProfileToRedis(profile)
	if err != nil {
		utils.Log("error on save profile when edit: %s", err.Error())
		SendError(bot, message.Chat.ID, errors.New(utils.ErrGotSomeProblems))
		return
	}
	SendText(bot, message.Chat.ID, "Данные о вашем профиле очищены")
}

func HandleProfile(bot *tgbotapi.BotAPI, message *tgbotapi.Message, profile *objects.Profile) {
	msgText := profile.FormatProfileMessage()
	msg := tgbotapi.NewMessage(profile.ChatID, msgText)
	msg.ReplyMarkup = profile.GetKeyboard()
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

	SendText(bot, chatID, "Ваш профиль обновлен")
	msgText := profile.FormatProfileMessage()
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = profile.GetKeyboard()
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
	if profile.Quote == 0 {
		SendText(bot, message.Chat.ID, "У вас закончилась квота на запросы. Пополните ее в разделе /payment")
		return
	}
	var messages []ai.Message
	messages = append(messages, ai.Message{Role: ai.RoleSystem, Content: utils.SystemPrompt})
	profileStr, err := profile.ProfileAIMessage()
	utils.Log("Profile message: %s", profileStr)
	if err != nil {
		utils.Log("Err formatting profile: %s", err.Error())
		SendError(bot, message.Chat.ID, errors.New(utils.ErrFillRequired))
		return
	}
	messages = append(messages, ai.Message{Role: ai.RoleUser, Content: profileStr})
	SendText(bot, message.Chat.ID, "Ожидаю нумерологический прогноз...")
	msgText, err := ai.GPTClient.SendMessage(messages)
	if err != nil {
		utils.Log("Err getting ai response: %s", err.Error())
		SendError(bot, message.Chat.ID, errors.New(utils.ErrGotSomeProblems))
		return
	}
	utils.Log("AI Answer: %s", msgText)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ParseMode = "Markdown"
	if SendMessage(bot, &msg) {
		profile.Quote -= 1
		err = database.SaveProfileToRedis(profile)
		if err != nil {
			utils.Log("error on save profile when edit: %s", err.Error())
			SendError(bot, message.Chat.ID, errors.New(utils.ErrGotSomeProblems))
			return
		}
	} else {
		utils.Log("error on send message. useless query")
	}
}

func HandleStop(bot *tgbotapi.BotAPI, message *tgbotapi.Message, profile *objects.Profile) {
	if profile.EditingField == "" {
		msgText := "Вы не находитесь в режиме изменения профиля"
		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
		SendMessage(bot, &msg)
		return
	}
	profile.EditingField = ""
	err := database.SaveProfileToRedis(profile)
	if err != nil {
		utils.Log("error on save profile when edit: %s", err.Error())
		SendError(bot, message.Chat.ID, errors.New(utils.ErrGotSomeProblems))
		return
	}
	msgText := "Ввод отменен:\n" + profile.FormatProfileMessage()
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyMarkup = profile.GetKeyboard()
	msg.ParseMode = "Markdown"
	SendMessage(bot, &msg)
}
