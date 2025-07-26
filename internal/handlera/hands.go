package handlera

import (
	"fmt"
	"net/http"

	"emobile/internal/dbase"
)

// DBPinger - Пинг базы данных
// router.HandleFunc("/ping", handlera.DBPinger).Methods("GET")
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
