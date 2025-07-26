package handlera

import (
	"emobile/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/google/uuid"
)

func (suite *TstHand) Test_01AddSub() {

	sub := models.Subscription{
		Service_name: "Yandex Plus",
		Price:        400,
		User_id:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		Start_date:   "01-02-2025",
		End_date:     "11-2025",
	}

	tests := []struct {
		name string
		//	dbEndPoint string
		sub    models.Subscription
		status int
		reply  string
	}{
		{
			name:   "Normaldu",
			sub:    sub,
			status: http.StatusOK,
			reply:  `{"status":"OK"}`,
		},
		{
			name: "No Price",
			sub: func() models.Subscription {
				s := sub
				s.Price = 0
				s.User_id = uuid.NewString()
				return s
			}(),
			status: http.StatusBadRequest,
			reply:  `{"status":"no price"}`,
		},
		{
			name: "No service name",
			sub: func() models.Subscription {
				s := sub
				s.Service_name = ""
				s.User_id = uuid.NewString()
				return s
			}(),
			status: http.StatusBadRequest,
			reply:  `{"status":"no Service_name"}`,
		},
		{
			name: "Bad start date",
			sub: func() models.Subscription {
				s := sub
				s.Start_date = "01-13-2022"
				s.User_id = uuid.NewString()
				s.End_date = ""
				return s
			}(),
			status: http.StatusBadRequest,
			reply:  `{"status":"bad START date"}`,
		},
		{
			name: "Nice start date, year 2 digits",
			sub: func() models.Subscription {
				s := sub
				s.Start_date = "01-10-22"
				s.User_id = uuid.NewString()
				s.End_date = ""
				return s
			}(),
			status: http.StatusOK,
			reply:  `{"status":"OK"}`,
		},
		{
			name: "End before start",
			sub: func() models.Subscription {
				s := sub
				s.Start_date = "24-02-2022"
				s.End_date = "08-08-2008"
				s.User_id = uuid.NewString()
				return s
			}(),
			status: http.StatusBadRequest,
			reply:  `{"status":"END date before START date"}`,
		},
		{
			name: "Wrong UUID",
			sub: func() models.Subscription {
				s := sub
				s.User_id = uuid.New().String() + "w"
				return s
			}(),
			status: http.StatusBadRequest,
			reply:  `{"status":"wrong uuid"}`,
		},
	}

	for _, tt := range tests {

		suite.Run(tt.name, func() {

			subM, err := json.Marshal(tt.sub)
			suite.Require().NoError(err)

			requestBody := strings.NewReader(string(subM))

			// Создание POST-запроса
			request := httptest.NewRequest(http.MethodPost, "/add", requestBody)

			// Установка заголовков
			request.Header.Set("Content-Type", "application/json")

			// Создание ResponseRecorder
			response := httptest.NewRecorder()
			// вызов хандлера
			CreateSub(response, request)

			res := response.Result()
			defer res.Body.Close()

			// Assert чтобы выполнилось сравнение tt.reply, string(resBody)
			suite.Assert().Equal(tt.status, res.StatusCode)

			// if tt.status != res.StatusCode {
			// 	resBody, err := io.ReadAll(res.Body)
			// 	suite.Require().NoError(err)
			// 	suite.Require().JSONEq(tt.reply, string(resBody))

			// }

		})
	}
}
