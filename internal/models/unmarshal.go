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
	sub.Sdt, err = parseDate(temp.StartDate)
	if err != nil {
		return err
	}
	sub.Edt, err = parseDate(temp.EndDate)
	if err != nil {
		return err
	}
	sub.Start_date = temp.StartDate
	sub.End_date = temp.EndDate

	return
}

func parseDate(date string) (t time.Time, err error) {

	// если дата пустая, возвращаем пустое (начальное) время, которое .IsZero() true
	if date == "" {
		return time.Time{}, nil
	}

	// парсим по день-месяц-год
	t, err = time.Parse("02-01-2006", date)
	// если ок - возвращаем
	if err == nil {
		return
	}
	// пробуем месяц-год
	t, err = time.Parse("01-2006", date)
	if err == nil {
		return
	}
	// парсим по день-месяц-год
	t, err = time.Parse("02-01-06", date)
	// если ок - возвращаем
	if err == nil {
		return
	}
	// пробуем месяц-год
	t, err = time.Parse("01-06", date)
	// if err == nil {
	// 	return
	// }

	return
}
