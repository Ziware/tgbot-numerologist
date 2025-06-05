package objects

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"tgbot-numerologist/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Profile struct {
	Username     string    `type:"internal" json:"username"`
	ChatID       int64     `type:"internal" json:"chat_id"`
	Predictions  int64     `type:"internal" json:"predictions"`
	Quote        int64     `type:"internal" json:"quote"`
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
	return Profile{Username: username, ChatID: chatId, Quote: 3}
}

func ParseDate(birthdate string) (time.Time, error) {
	parsed, err := time.Parse("02.01.2006", birthdate)
	if err != nil {
		utils.Log("error while parse date %s: %s", birthdate, err.Error())
		return parsed, errors.New(utils.ErrWrongTimeFormat)
	}
	return parsed, nil
}

func (p *Profile) ResetProfile() {
	p.Name = ""
	p.Surname = ""
	p.BirthDate = time.Time{}
	p.Bio = ""
	p.WorkPlace = ""
	p.StudyPlace = ""
	p.Hobby = ""
}

func (p *Profile) FormatProfileMessage() string {
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
	res += formatRow("Имя", p.Name, "Name")
	res += formatRow("Фамилия", p.Surname, "Surname")
	if p.BirthDate.IsZero() {
		res += formatRow("Дата рождения", "", "BirthDate")
	} else {
		res += formatRow("Дата рождения", p.BirthDate.Format("02.01.2006"), "BirthDate")
	}
	res += formatRow("Место работы", p.WorkPlace, "WorkPlace")
	res += formatRow("Место учёбы", p.StudyPlace, "StudyPlace")
	res += formatRow("Хобби", p.Hobby, "Hobby")
	res += formatRow("Биография", p.Bio, "Bio")
	return res
}

func (p *Profile) GetKeyboard() tgbotapi.InlineKeyboardMarkup {
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

func (p *Profile) GetPaymentKeyboard() tgbotapi.InlineKeyboardMarkup {
	buttons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Увеличить квоту", "pay"),
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, btn := range buttons {
		row := tgbotapi.NewInlineKeyboardRow(btn)
		keyboard = append(keyboard, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

func (p *Profile) ProfileAIMessage() (string, error) {
	v := reflect.ValueOf(*p)
	t := v.Type()

	var res string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		fieldType, ok := field.Tag.Lookup("type")

		if !ok {
			continue
		}

		var strValue string
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			if tt, ok := fieldValue.Interface().(time.Time); ok && !tt.IsZero() {
				strValue = tt.Format("02.01.2006")
			} else {
				strValue = ""
			}
		} else {
			strValue = fmt.Sprint(fieldValue.Interface())
		}

		switch fieldType {
		case "internal":
			continue
		case "required":
			if len(strValue) == 0 {
				return "", fmt.Errorf("required field '%s' is empty", field.Name)
			}
			res += field.Name + ": " + strValue + "\n"
		case "optional":
			if len(strValue) != 0 {
				res += field.Name + ": " + strValue + "\n"
			}
		}
	}
	return res, nil
}
