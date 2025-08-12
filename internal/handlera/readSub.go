package handlera

import (
	"emobile/internal/models"
	"encoding/json"
	"errors"
	"net/http"
)

// ReadSub godoc
// @Summary Получить подписки
// @Description Возвращает список подписок по заданным параметрам
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Параметры поиска (обязательно service_name и user_id)"
// @Success 200 {array} models.Subscription
// @Failure 400 {object} string "Неверный формат запроса или отсутствуют обязательные поля"
// @Failure 500 {object} string "Ошибка сервера"
// @Router /read [post]
func (db *DBstruct) ReadSub(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "application/json")

	readSub := models.Subscription{}
	err := json.NewDecoder(req.Body).Decode(&readSub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusBadRequest)
		return
	}

	// в запросе read обязятельныц поля Service_name и User_id
	if readSub.Service_name == "" {
		http.Error(rwr, "no service name", http.StatusBadRequest)
		return
	}
	if readSub.User_id == "" {
		http.Error(rwr, "no user_id", http.StatusBadRequest)
		return
	}

	// models.Inter, err = dbase.NewPostgresPool(req.Context(), models.DSN)
	// if err != nil {
	// 	http.Error(rwr, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer models.Inter.CloseDB()

	// subs []models.Subscription
	subs, err := db.Inter.ReadSub(req.Context(), readSub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}

	rwr.WriteHeader(http.StatusOK)

	if len(subs) != 0 {
		models.Logger.Info("Найдено", "подписок", len(subs))
		json.NewEncoder(rwr).Encode(subs)
		return
	}

	models.Logger.Info("Read - Не найдено записей")
	ret := models.RetStruct{
		Name: "Не найдено записей, удовлетворяющих запросу",
		Cunt: 0,
	}
	json.NewEncoder(rwr).Encode(ret)

}

// UpdateSub godoc
// @Summary Обновление подписки
// @Description Обновляет данные подписки в базе данных
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Данные подписки для обновления (обязательные: service_name и user_id)"
// @Success 200 {object} models.RetStruct
// @Failure 400 {object} object "Неверный запрос"
// @Failure 500 {object} object "Ошибка сервера"
// @Router /update [put]
func (db *DBstruct) UpdateSub(rwr http.ResponseWriter, req *http.Request) {

	rwr.Header().Set("Content-Type", "application/json")

	readSub := models.Subscription{}
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

	// models.Inter, err = dbase.NewPostgresPool(req.Context(), models.DSN)
	// if err != nil {
	// 	rwr.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(rwr).Encode(err)
	// 	return
	// }
	// defer models.Inter.CloseDB()

	cTag, err := db.Inter.UpdateSub(req.Context(), readSub)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rwr).Encode(err)
		return
	}

	ret := models.RetStruct{
		Name: "Обновлено записей",
		Cunt: cTag.RowsAffected(),
	}
	if cTag.RowsAffected() == 0 {
		ret.Name = "Не найдено записей, удовлетворяющих запросу"
	}

	models.Logger.Info("UPDATE", "", ret)

	rwr.WriteHeader(http.StatusOK)
	json.NewEncoder(rwr).Encode(ret)
}
