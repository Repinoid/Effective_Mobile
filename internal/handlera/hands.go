package handlera

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"emobile/internal/dbase"
	"emobile/internal/models"

	"github.com/google/uuid"
)

// router.HandleFunc("/ping", handlera.DBPinger).Methods("GET")
// DBPinger - Пинг базы данных
func DBPinger(rwr http.ResponseWriter, req *http.Request) {

	err := dbase.Ping(req.Context())
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rwr, `{"Error":"%v"}`, err)
		return
	}
	rwr.WriteHeader(http.StatusOK)
	fmt.Fprintf(rwr, `{"status":"StatusOK"}`)
}

// CreateSub создаёт новую запись о подписке
func CreateSub(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")

	// telo, err := io.ReadAll(req.Body)
	// if err != nil {
	// 	rwr.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
	// 	return
	// }
	// defer req.Body.Close()

	// sub := models.Subscription{}
	// err = json.Unmarshal(telo, &sub)
	// if err != nil {
	// 	rwr.WriteHeader(http.StatusBadRequest) // с некорректным  значением возвращать http.StatusBadRequest.
	// 	json.NewEncoder(rwr).Encode(err)
	// 	return
	// }

	sub := models.Subscription{}
	err := json.NewDecoder(req.Body).Decode(&sub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusBadRequest)
		return
	}

	if sub.Service_name == "" {
		rwr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rwr).Encode(errors.New("no service name"))
		return
	}
	if sub.Price == 0 {
		rwr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rwr).Encode(errors.New("no price"))
		return
	}
	if sub.User_id == "" {
		rwr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rwr).Encode(errors.New("no user_id"))
		return
	}
	_, err = uuid.Parse(sub.User_id)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rwr).Encode(errors.New("bad user_id"))
		return
	}

	if sub.Start_date == "" {
		rwr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rwr).Encode(errors.New("no start date"))
		return
	}
	Start_date, err := parseDate(sub.Start_date)
	if err != nil {
		rwr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rwr).Encode(err)
		return
	}
	sub.Sdt = Start_date

	if sub.End_date != "" {
		End_date, err := parseDate(sub.End_date)
		if err != nil {
			rwr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rwr).Encode(err)
			return
		}
		if End_date.Before(Start_date) {
			rwr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rwr).Encode(errors.New("end date before start"))
			return

		}
		sub.Edt = End_date
	}
	db, err := dbase.NewPostgresPool(req.Context(), models.DSN)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rwr).Encode(err)
		return
	}
	defer db.DB.Close()

	err = db.AddSub(req.Context(), sub)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rwr).Encode(err)
		return
	}

	rwr.WriteHeader(http.StatusOK)

	// в поле Message err сообщение об ошибке типа
	// "Message": "duplicate key value violates unique constraint \"subscriptions_user_id_key\""
	// `{"Message":"OK"}` - для унификации, если парсить возврат из хандлера по полю "Message"
	fmt.Fprintf(rwr, `{"Message":"OK"}`)

}

func parseDate(date string) (t time.Time, err error) {
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
