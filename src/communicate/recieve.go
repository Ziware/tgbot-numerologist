package communicate

import (
	"errors"

	"tgbot-numerologist/database"
	"tgbot-numerologist/objects"
	"tgbot-numerologist/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetProfile(username string, chatID int64) (*objects.Profile, error) {
	exists, err := database.UserExistsInRedis(username)
	if err != nil {
		utils.Log("error check exists in redis: %v", err)
		return nil, errors.New(utils.ErrGotSomeProblems)
	}
	if exists {
		profile, err := database.GetProfileFromRedis(username)
		if err != nil {
			utils.Log("error get profile from redis: %v", err)
			return nil, errors.New(utils.ErrGotSomeProblems)
		}
		return profile, nil
	}
	utils.Log("profile for user %s not exists, create with chat id %d", username, chatID)
	profile := objects.NewProfile(username, chatID)
	err = database.SaveProfileToRedis(&profile)
	if err != nil {
		utils.Log("error save to redis: %v", err)
		return nil, errors.New(utils.ErrGotSomeProblems)
	}
	return &profile, nil
}

func StartReceivingUpdates(bot *tgbotapi.BotAPI) {
	utils.Log("Start Receiving Updates")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		var msg *tgbotapi.Message = nil
		var profile *objects.Profile = nil
		var err error
		if update.CallbackQuery != nil {
			msg = update.CallbackQuery.Message
			profile, err = GetProfile(update.CallbackQuery.From.UserName, msg.Chat.ID)
			utils.Log("Get profile: %s:%s", profile.Username, profile.EditingField)
			if err != nil {
				SendError(bot, msg.Chat.ID, errors.New(utils.ErrGotSomeProblems))
				continue
			}
		}
		if update.Message != nil {
			msg = update.Message
			LogMessage(bot, msg)
			profile, err = GetProfile(msg.From.UserName, msg.Chat.ID)
			utils.Log("Get profile: %s:%s", profile.Username, profile.EditingField)
			if err != nil {
				SendError(bot, msg.Chat.ID, errors.New(utils.ErrGotSomeProblems))
				continue
			}
		}
		if msg == nil || profile == nil {
			continue
		}
		if update.CallbackQuery != nil {
			DetermineCallback(bot, update.CallbackQuery, profile)
			continue
		}

		if msg.IsCommand() {
			DetermineCommand(bot, msg, profile)
			continue
		}

		if profile.EditingField != "" {
			HandleEditMessage(bot, msg, profile)
			continue
		}

		SendCommon(bot, msg)
	}
}

func DetermineCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, profile *objects.Profile) {
	switch message.Command() {
	case "start":
		HandleIntro(bot, message)
	case "intro":
		HandleIntro(bot, message)
	case "help":
		HandleHelp(bot, message)
	case "feedback":
		HandleFeedback(bot, message)
	case "payment":
		HandlePayment(bot, message, profile)
	case "reset":
		HandleReset(bot, message, profile)
	case "profile":
		HandleProfile(bot, message, profile)
	case "predictions":
		HandlePredictions(bot, message, profile)
	case "stop":
		HandleStop(bot, message, profile)
	default:
		SendText(bot, message.Chat.ID, utils.ErrUnknownCommand)
	}
}

func DetermineCallback(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, profile *objects.Profile) {
	chatID := callbackQuery.Message.Chat.ID
	switch callbackQuery.Data {
	case "pay":
		HandlePayButton(bot, callbackQuery, profile)
		return
	}
	SendText(bot, chatID, "Напишите /stop для отмены ввода")
	switch callbackQuery.Data {
	case "edit_name":
		SendText(bot, chatID, "Введите ваше имя:")
	case "edit_surname":
		SendText(bot, chatID, "Введите вашу фамилию:")
	case "edit_birthdate":
		SendText(bot, chatID, "Введите вашу дату рождения в формате dd.mm.yyyy:")
	case "edit_bio":
		SendText(bot, chatID, "Введите вашу биографию:")
	case "edit_workplace":
		SendText(bot, chatID, "Введите ваше место работы:")
	case "edit_studyplace":
		SendText(bot, chatID, "Введите ваше место учёбы:")
	case "edit_hobby":
		SendText(bot, chatID, "Опишите ваше хобби:")
	}
	profile.EditingField = callbackQuery.Data
	utils.Log("enter edit phase: %s", profile.EditingField)

	err := database.SaveProfileToRedis(profile)
	if err != nil {
		utils.Log("error on save profile when edit: %s", err.Error())
		SendError(bot, callbackQuery.Message.Chat.ID, errors.New(utils.ErrGotSomeProblems))
		return
	}
	utils.Log("successfully saved data to redis")

	p, err := database.GetProfileFromRedis(profile.Username)
	utils.Log("get already written: %s:%s", p.Username, p.EditingField)

	bot.Request(tgbotapi.NewCallback(callbackQuery.ID, "Ожидаю ввода..."))
}
