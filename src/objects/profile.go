package objects

import (
	"errors"
	"time"

	"tgbot-numerologist/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Profile struct {
	Username     string    `type:"internal" json:"username"`
	ChatID       int64     `type:"internal" json:"chat_id"`
	EditingField string    `type:"internal" json:"editing_field"`
	Name         string    `type:"required" json:"name"`
	Surname      string    `type:"optional" json:"surname"`
	BirthDate    time.Time `type:"required" json:"birth_date"`
	Bio          string    `type:"optional" json:"bio"`
	WorkPlace    string    `type:"optional" json:"work_place"`
	StudyPlace   string    `type:"optional" json:"study_place"`
	Hobby        string    `type:"optional" json:"hobby"`
}

func NewProfile(username string, chatId int64) Profile {
	return Profile{Username: username, ChatID: chatId}
}

func ParseDate(birthdate string) (time.Time, error) {
	parsed, err := time.Parse("02.01.2006", birthdate)
	if err != nil {
		utils.Log("error while parse date %s: %s", birthdate, err.Error())
		return parsed, errors.New(utils.ErrWrongTimeFormat)
	}
	return parsed, nil
}

func FormatProfileMessage(profile *Profile) string {
	requiredFields := []string{"Name", "BirthDate"}
	fieldStatus := func(field string) string {
		for _, requiredField := range requiredFields {
			if field == requiredField {
				return "_(обязательно)_"
			}
		}
		return ""
	}
	formatRow := func(key, value, field string) string {
		if value == "" {
			value = "-"
		}
		res := key + ": " + value
		status := fieldStatus(field)
		if len(status) != 0 && value == "-" {
			res += " " + status
		}
		res += "\n"
		return res
	}
	res := "*Ваш профиль*:\n"
	res += formatRow("Имя", profile.Name, "Name")
	res += formatRow("Фамилия", profile.Surname, "Surname")
	if profile.BirthDate.IsZero() {
		res += formatRow("Дата рождения", "", "BirthDate")
	} else {
		res += formatRow("Дата рождения", profile.BirthDate.Format("02.01.2006"), "BirthDate")
	}
	res += formatRow("Место работы", profile.WorkPlace, "WorkPlace")
	res += formatRow("Место учёбы", profile.StudyPlace, "StudyPlace")
	res += formatRow("Хобби", profile.Hobby, "Hobby")
	res += formatRow("Биография", profile.Bio, "Bio")
	return res
}

func ProfileKeyboard() tgbotapi.InlineKeyboardMarkup {
	buttons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Изменить имя", "edit_name"),
		tgbotapi.NewInlineKeyboardButtonData("Изменить фамилию", "edit_surname"),
		tgbotapi.NewInlineKeyboardButtonData("Изменить дату рождения", "edit_birthdate"),
		tgbotapi.NewInlineKeyboardButtonData("Изменить место работы", "edit_workplace"),
		tgbotapi.NewInlineKeyboardButtonData("Изменить место учёбы", "edit_studyplace"),
		tgbotapi.NewInlineKeyboardButtonData("Изменить хобби", "edit_hobby"),
		tgbotapi.NewInlineKeyboardButtonData("Изменить биографию", "edit_bio"),
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, btn := range buttons {
		row := tgbotapi.NewInlineKeyboardRow(btn)
		keyboard = append(keyboard, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}
