package handlera

import (
	"emobile/internal/dbase"
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
func SumSub(rwr http.ResponseWriter, req *http.Request) {

	rwr.Header().Set("Content-Type", "application/json")

	readSub := models.Subscription{}
	err := json.NewDecoder(req.Body).Decode(&readSub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusBadRequest)
		return
	}

	// должны быть определены Service_name или User_id, и диапазон дат
	if (readSub.Service_name == "" && readSub.User_id == "") ||
		readSub.End_date == nil || readSub.Start_date == nil {
		http.Error(rwr, "не все данные указаны", http.StatusBadRequest)
		return
	}

	models.Inter, err = dbase.NewPostgresPool(req.Context(), models.DSN)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	defer models.Inter.CloseDB()

	summa, err := models.Inter.SumSub(req.Context(), readSub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}

	rwr.WriteHeader(http.StatusOK)

	ret := models.RetStruct{
		Name: "Сумма подписок",
		Cunt: summa,
	}

	if summa == 0 {
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
func DeleteSub(rwr http.ResponseWriter, req *http.Request) {

	rwr.Header().Set("Content-Type", "application/json")

	readSub := models.Subscription{}
	err := json.NewDecoder(req.Body).Decode(&readSub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusBadRequest)
		return
	}

	models.Inter, err = dbase.NewPostgresPool(req.Context(), models.DSN)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	defer models.Inter.CloseDB()

	cTag, err := models.Inter.DeleteSub(req.Context(), readSub)
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
