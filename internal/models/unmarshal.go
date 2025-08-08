package models

import (
	"encoding/json"
	"errors"
	"time"
)
// переопределение метода анмаршаллинга для типа Subscription
func (sub *Subscription) UnmarshalJSON(data []byte) (err error) {

	// Создаем временный тип, чтобы избежать рекурсии при Unmarshal
	type Alias Subscription
	temp := &struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(sub),
	}
	// Разбираем JSON во временную структуру
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	// дата из строки - в time.Time
	sub.Start_date, err = ParseDate(temp.StartDate)
	if err != nil {
		return err
	}
	sub.End_date, err = ParseDate(temp.EndDate)
	if err != nil {
		return err
	}

	return
}

// 	ParseDate принимает строковую дату, возвращает time.Time или nil. Отсекает день месяца, устанавливает в 1е число месяца - подписка помесячно, даты не важны
func ParseDate(date string) (tim any, err error) {

	// если дата пустая
	if date == "" {
		return nil, nil
	}

	// парсим по день-месяц-год
	// если ок - возвращаем
	if t, err := time.Parse("02-01-2006", date); err == nil {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()), nil
	}
	// пробуем месяц-год
	if t, err := time.Parse("01-2006", date); err == nil {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()), nil
	}
	// парсим по день-месяц-год
	if t, err := time.Parse("02-01-06", date); err == nil {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()), nil
	}
	// пробуем месяц-год
	if t, err := time.Parse("01-06", date); err == nil {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()), nil
	}
	if t, err := time.Parse(time.RFC3339, date); err == nil {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()), nil
	}

	return nil, errors.New("неверный формат даты")
}

// MakeTT используется в функциях тестов, преобразует поля со строковыми датами в time.Time
func MakeTT(sub *Subscription) (err error) {

	switch sub.Start_date.(type) {
	case string:
		sub.Start_date, _ = ParseDate(sub.Start_date.(string))
	}
	switch sub.End_date.(type) {
	case string:
		sub.End_date, _ = ParseDate(sub.End_date.(string))
	}
	return
}
