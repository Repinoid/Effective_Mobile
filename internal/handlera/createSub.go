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
// @Summary Database health check
// @Description Checks if database connection is alive
// @Produce json
// @Success 200 {object} map[string]string "Database is reachable"
// @Failure 500 {object} map[string]string "Database connection error"
// @Router / [get]
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

// CreateSub godoc
// @Summary Create a new subscription
// @Description Add a new subscription to the database
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Subscription data"
// @Success 200 {object} models.RetStruct
// @Failure 400 {object} object "Validation error"
// @Failure 500 {object} object "Internal server error"
// @Router /add [post]
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

	models.Inter, err = dbase.NewPostgresPool(req.Context(), models.DSN)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rwr).Encode(err)
		return
	}
	defer models.Inter.CloseDB()

	cTag, err := models.Inter.AddSub(req.Context(), sub)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rwr).Encode(err)
		return
	}

	ret := models.RetStruct{
		Name: "Внесено записей",
		Cunt: cTag.RowsAffected(),
	}

	rwr.WriteHeader(http.StatusOK)
	json.NewEncoder(rwr).Encode(ret)

}
