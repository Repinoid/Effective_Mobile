package handlera

import (
	"emobile/internal/models"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (suite *TstHand) Test_03UpdateSub() {
	// update запись по запросу subForUpdate
	subForUpdate := models.ReadSubscription{
		Service_name: "Yandex Plus", //
		Price:        666,
		User_id:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		Start_date:   "08-08-08",
		End_date:     "24-02-22",
	}

	subM, err := json.Marshal(subForUpdate)
	suite.Require().NoError(err)

	requestBody := strings.NewReader(string(subM))

	request := httptest.NewRequest(http.MethodPut, "/update", requestBody)

	// Создание ResponseRecorder
	response := httptest.NewRecorder()
	// вызов хандлера
	UpdateSub(response, request)

	res := response.Result()
	defer res.Body.Close()

	// HTTP put UPDATE должен вернуть http.StatusOK
	suite.Require().Equal(http.StatusOK, res.StatusCode)

	// Составляем запрос READ для чтения только что UPDATEd записи
	subReadUpdated := models.ReadSubscription{
		Service_name: "Yandex Plus",
		User_id:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
	}

	subM, err = json.Marshal(subReadUpdated)
	suite.Require().NoError(err)

	requestBody = strings.NewReader(string(subM))

	// Создание HTTP POST-запроса на чтение записи
	request = httptest.NewRequest(http.MethodPost, "/read", requestBody)

	// Установка заголовков
	request.Header.Set("Content-Type", "application/json")

	// Создание ResponseRecorder
	response = httptest.NewRecorder()
	// вызов хандлера
	ReadSub(response, request)

	res = response.Result()
	defer res.Body.Close()

	// http.StatusOK should be
	suite.Require().Equal(http.StatusOK, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	suite.Require().NoError(err)

	// размаршалливаем список
	subs := []models.ReadSubscription{}
	err = json.Unmarshal(resBody, &subs)
	suite.Require().NoError(err)
	// должна быть всего одна запись
	suite.Require().Equal(1, len(subs))

	// убеждаемся, что запись обновилась
	suite.Require().Equal(subForUpdate.Service_name, subs[0].Service_name)
	suite.Require().Equal(subForUpdate.User_id, subs[0].User_id)
	suite.Require().Equal(subForUpdate.Price, subs[0].Price)
	suite.Require().Equal(subForUpdate.Sdt, subs[0].Sdt)
	suite.Require().Equal(subForUpdate.Edt, subs[0].Edt)
}
