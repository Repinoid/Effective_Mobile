package handlera

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
	// если при непустой конечной дате она раньше начальной
	if sub.End_date != "" && sub.Edt.Before(sub.Sdt) {
		rwr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rwr).Encode(errors.New("end date before start"))
		return
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
