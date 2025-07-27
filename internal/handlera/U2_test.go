package handlera

import (
	"emobile/internal/models"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// type errMessage struct {
// 	Message   string `json:"Message"`
// 	Detail    string `json:"Detail"`
// 	TableName string `json:"TableName"`
// }

func (suite *TstHand) Test_02ReadSub() {

	// Völkischer Beobachter   Avanti

	sub := models.ReadSubscription{
		Service_name: "Yandex Plus",
		//	Price:        400,
		User_id: "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		//	Start_date:   "01-02-2025",
		//	End_date:     "11-2025",
	}

	tests := []struct {
		name string
		//	dbEndPoint string
		sub    models.ReadSubscription
		status int
		reply  string
	}{
		{
			name:   "Normaldu",
			sub:    sub,
			status: http.StatusOK,
			reply:  `{"status":"OK"}`,
		},
	}

	for _, tt := range tests {

		suite.Run(tt.name, func() {

			subM, err := json.Marshal(tt.sub)
			suite.Require().NoError(err)

			requestBody := strings.NewReader(string(subM))

			// Создание POST-запроса
			request := httptest.NewRequest(http.MethodPost, "/read", requestBody)

			// Установка заголовков
			request.Header.Set("Content-Type", "application/json")

			// Создание ResponseRecorder
			response := httptest.NewRecorder()
			// вызов хандлера
			ReadSub(response, request)

			res := response.Result()
			defer res.Body.Close()

			// Assert чтобы выполнилось сравнение tt.reply, string(resBody)
			suite.Require().Equal(tt.status, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)
			suite.Require().NoError(err)

			// размаршалливаем список подписок
			subs := []models.ReadSubscription{}
			err = json.Unmarshal(resBody, &subs)
			suite.Require().NoError(err)
			// должно быть 2 записи
			suite.Require().Equal(1, len(subs))
			// сравниваем Service_name и User_id первой записи
			suite.Require().Equal(sub.Service_name, subs[0].Service_name)
			suite.Require().Equal(sub.User_id, subs[0].User_id)
		})
	}

}
