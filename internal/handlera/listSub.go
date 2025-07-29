package handlera

import (
	"encoding/json"
	"net/http"

	"emobile/internal/dbase"
	"emobile/internal/models"
)

// ListSub godoc
// @Summary Получить список всех подписок
// @Description Возвращает полный список всех подписок из базы данных
// @Tags Подписки
// @Produce json
// @Success 200 {array} models.Subscription
// @Failure 500 {object} object "Ошибка сервера"
// @Router /list [get]
func ListSub(rwr http.ResponseWriter, req *http.Request) {

	db, err := dbase.NewPostgresPool(req.Context(), models.DSN)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.DB.Close()

	// запрос в БД на получения списка всех подписок
	subs, err := db.ListSub(req.Context())
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}

	rwr.WriteHeader(http.StatusOK)

	if len(subs) != 0 {
		json.NewEncoder(rwr).Encode(subs)
	} else {
		ret := models.RetStruct{
			Name: "Нет записей в подписках",
			Cunt: 0,
		}
		json.NewEncoder(rwr).Encode(ret)
	}

}
