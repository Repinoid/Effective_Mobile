package handlera

import (
	"encoding/json"
	"net/http"

	"emobile/internal/dbase"
	"emobile/internal/models"
)

func ListSub(rwr http.ResponseWriter, req *http.Request) {

	db, err := dbase.NewPostgresPool(req.Context(), models.DSN)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rwr).Encode(err)
		return
	}
	defer db.DB.Close()

	// запрос в БД на получения списка всех подписок 
	subs, err := db.ListSub(req.Context())
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rwr).Encode(err)
		return
	}

	rwr.WriteHeader(http.StatusOK)

	// возврат из хандлера маршалленого списка подписок
	// Encode writes the JSON encoding of v to the stream,
	// with insignificant space characters elided, followed by a newline character.
	json.NewEncoder(rwr).Encode(subs)

}
