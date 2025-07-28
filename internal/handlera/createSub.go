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

// DBPinger godoc
// @Summary Проверка соединения с БД
// @Description Проверяет доступность и работоспособность базы данных
// @Tags health
// @Produce json
// @Success 200 {object} DBStatusResponse
// @Failure 500 {object} ErrorResponse
// @Router / [get]
func DBPinger(rwr http.ResponseWriter, req *http.Request) {

	err := dbase.Ping(req.Context())
	if err != nil {
		http.Error(rwr, "Нет соединения с сервером", http.StatusInternalServerError)
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

	cTag, err := db.AddSub(req.Context(), sub)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rwr).Encode(err)
		return
	}

	rowsAffected := cTag.RowsAffected()
	ret := struct {
		Name string
		rows int64
	}{"Внесено записей", rowsAffected}

	json.NewEncoder(rwr).Encode(ret)

	rwr.WriteHeader(http.StatusOK)

}
