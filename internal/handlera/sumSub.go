package handlera

import (
	"database/sql"
	"emobile/internal/models"
	"encoding/json"
	"net/http"
)

// SumSub godoc
// @Summary Расчет суммы подписок
// @Description Возвращает сумму подписок по заданным параметрам
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Параметры для расчета суммы"
// @Success 200 {object} models.RetStruct
// @Failure 400 {object} string "Неверный формат запроса или отсутствуют обязательные поля"
// @Failure 500 {object} string "Ошибка сервера"
// @Router /summa [post]
func (db *InterStruct) SumHandler(rwr http.ResponseWriter, req *http.Request) {

	rwr.Header().Set("Content-Type", "application/json")

	readSub := models.Subscription{}
	err := json.NewDecoder(req.Body).Decode(&readSub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusBadRequest)
		return
	}

	if (readSub.Service_name == "" && readSub.User_id == "") ||
		readSub.Edt.IsZero() || readSub.Sdt.IsZero() {
		http.Error(rwr, "не все данные указаны", http.StatusBadRequest)
		return
	}

	// models.Inter, err = dbase.NewPostgresPool(req.Context(), models.DSN)
	// if err != nil {
	// 	http.Error(rwr, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer models.Inter.CloseDB()

	summa, err := db.Inter.SumSub(req.Context(), readSub)
	if err != nil && err != sql.ErrNoRows {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}

	rwr.WriteHeader(http.StatusOK)

	ret := models.RetStruct{
		Name: "Сумма подписок",
		Cunt: summa,
	}

	if summa == 0 || err == sql.ErrNoRows {
		ret.Name = "Нет таких подписок"
	}
	models.Logger.Info("Сумма подписок ", "", ret)

	json.NewEncoder(rwr).Encode(ret)
}

// DeleteSub godoc
// @Summary Удаление подписки
// @Description Удаляет подписку по переданным данным
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Данные подписки для удаления"
// @Success 200 {object} models.RetStruct
// @Failure 400 {object} string "Неверный формат запроса"
// @Failure 500 {object} string "Ошибка сервера"
// @Router /delete [delete]
func (db *InterStruct) DeleteHandler(rwr http.ResponseWriter, req *http.Request) {

	rwr.Header().Set("Content-Type", "application/json")

	readSub := models.Subscription{}
	err := json.NewDecoder(req.Body).Decode(&readSub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusBadRequest)
		return
	}

	// models.Inter, err = dbase.NewPostgresPool(req.Context(), models.DSN)
	// if err != nil {
	// 	http.Error(rwr, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer models.Inter.CloseDB()

	cTag, err := db.Inter.DeleteSub(req.Context(), readSub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}

	rwr.WriteHeader(http.StatusOK)

	ret := models.RetStruct{
		Name: "Удалено записей",
		Cunt: cTag.RowsAffected(),
	}

	if cTag.RowsAffected() == 0 {
		ret.Name = "Не найдено записей на удаление, удовлетворяющих запросу"
	}

	models.Logger.Info("DELETE ", "", ret)

	json.NewEncoder(rwr).Encode(ret)
}
