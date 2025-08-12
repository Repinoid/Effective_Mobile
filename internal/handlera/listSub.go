package handlera

import (
	"encoding/json"
	"net/http"
	"strconv"

	"emobile/internal/config"
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
func (db *InterStruct) ListSub(rwr http.ResponseWriter, req *http.Request) {

	// Получение параметров страницы
	pageStr := req.URL.Query().Get("page")
	pageSizeStr := req.URL.Query().Get("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = config.Configuration.PageSize
	}

	offset := (page - 1) * pageSize

	// models.Inter, err = dbase.NewPostgresPool(req.Context(), models.DSN)
	// if err != nil {
	// 	http.Error(rwr, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer models.Inter.CloseDB()

	// запрос в БД на получения списка всех подписок
	subs, err := db.Inter.ListSub(req.Context(), pageSize, offset)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}

	rwr.WriteHeader(http.StatusOK)

	if len(subs) != 0 {
		models.Logger.Info("Cписок", "подписки", subs)
		json.NewEncoder(rwr).Encode(subs)
		return
	}
	models.Logger.Info("Нет записей в подписках")

	ret := models.RetStruct{
		Name: "Нет записей в подписках",
		Cunt: 0,
	}
	json.NewEncoder(rwr).Encode(ret)

}
