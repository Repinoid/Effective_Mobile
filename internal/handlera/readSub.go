package handlera

import (
	"emobile/internal/dbase"
	"emobile/internal/models"
	"encoding/json"
	"errors"
	"net/http"
)

func ReadSub(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")

	// telo, err := io.ReadAll(req.Body)
	// if err != nil {
	// 	rwr.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(rwr, `{"status":"StatusBadRequest"}`)
	// 	return
	// }
	// defer req.Body.Close()

	// readSub := models.ReadSubscription{}
	// err = json.Unmarshal(telo, &readSub)
	// if err != nil {
	// 	rwr.WriteHeader(http.StatusBadRequest) // с некорректным  значением возвращать http.StatusBadRequest.
	// 	json.NewEncoder(rwr).Encode(err)
	// 	return
	// }

	readSub := models.ReadSubscription{}
	err := json.NewDecoder(req.Body).Decode(&readSub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusBadRequest)
		return
	}

	// в запросе read обязятельныц поля Service_name и User_id
	if readSub.Service_name == "" {
		rwr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rwr).Encode(errors.New("no service name"))
		return
	}
	if readSub.User_id == "" {
		rwr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rwr).Encode(errors.New("no user_id"))
		return
	}

	// если присутствуют даты, конвертируем их в timastamp и заполняем .Sdt .Еdt
	if readSub.Start_date != "" {
		Start_date, err := parseDate(readSub.Start_date)
		if err != nil {
			rwr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rwr).Encode(err)
			return
		}
		readSub.Sdt = Start_date
	}
	if readSub.End_date != "" {
		End_date, err := parseDate(readSub.End_date)
		if err != nil {
			rwr.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rwr).Encode(err)
			return
		}
		readSub.Edt = End_date
	}

	db, err := dbase.NewPostgresPool(req.Context(), models.DSN)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rwr).Encode(err)
		return
	}
	defer db.DB.Close()

	subs, err := db.ReadSub(req.Context(), readSub)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rwr).Encode(err)
		return
	}
	_ = subs

	rwr.WriteHeader(http.StatusOK)

	json.NewEncoder(rwr).Encode(subs)

}
