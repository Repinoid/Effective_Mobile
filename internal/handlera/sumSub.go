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
		readSub.Edt.IsZero() || readSub.Sdt.IsZero() {
		http.Error(rwr, "не все данные указаны", http.StatusBadRequest)
		return
	}

	db, err := dbase.NewPostgresPool(req.Context(), models.DSN)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.DB.Close()

	summa, err := db.SumSub(req.Context(), readSub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}

	rwr.WriteHeader(http.StatusOK)

	ret := models.RetStruct{
		Name: "Сумма подписок",
		Cunt: summa,
	}
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

	db, err := dbase.NewPostgresPool(req.Context(), models.DSN)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.DB.Close()

	cTag, err := db.DeleteSub(req.Context(), readSub)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}

	rwr.WriteHeader(http.StatusOK)

	ret := models.RetStruct{
		Name: "Удалено записей",
		Cunt: cTag.RowsAffected(),
	}

	json.NewEncoder(rwr).Encode(ret)
}
