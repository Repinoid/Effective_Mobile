package models

import (
	"encoding/json"
	"time"
)

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
	sub.Sdt, err = ParseDate(temp.StartDate)
	if err != nil {
		return err
	}
	// if sub.Sdt.(time.Time).IsZero() {
	// 	sub.Sdt = nil
	// }
	sub.Edt, err = ParseDate(temp.EndDate)
	if err != nil {
		return err
	}
	// if sub.Edt.(time.Time).IsZero() {
	// 	sub.Edt = nil
	// }

	sub.Start_date = temp.StartDate
	sub.End_date = temp.EndDate

	return
}

// возвращаем 1е число месяца - подписка помесячно, даты не важны
func ParseDate(date string) (tim any, err error) {

	// если дата пустая, возвращаем пустое (начальное) время, которое .IsZero() true
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

	return nil, err
}

func MakeTT(sub *Subscription) (err error) {
	sub.Sdt, err = ParseDate(sub.Start_date)
	if err != nil {
		return
	}
	sub.Edt, err = ParseDate(sub.End_date)

	// if sub.Edt.(time.Time).IsZero() {
	// 	sub.Edt = nil
	// }
	// if sub.Sdt.(time.Time).IsZero() {
	// 	sub.Sdt = nil
	// }

	// if sub.Edt.IsZero() {
	// 	sub.Edt = time.Date(9999, time.December, 31, 23, 59, 59, 999999999, time.UTC)
	// }
	return

}
