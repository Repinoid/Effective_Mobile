package main

import (
	"emobile/internal/models"
	"net/http"

	"github.com/go-resty/resty/v2"
)

var sub = models.Subscription{
	Service_name: "Жейминь Жибао1",
	Price:        400,
	User_id:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
	Start_date:   "01-02-2020",
	End_date:     "11-2029",
}

func (suite *TS) Test_01() {

	httpc := resty.New().SetBaseURL("http://localhost:8080")

	req := httpc.R().SetHeader("Content-Type", "application/json").SetDoNotParseResponse(false).
		SetBody(sub)

	resp, err := req.Post("/add")
	suite.Require().NoError(err, "req.Post add)")

	suite.Require().Equal(http.StatusOK, resp.StatusCode())

	suite.Require().JSONEq(`{"Cunt":1, "Name":"Внесено записей"}`, resp.String())

}

func (suite *TS) TestExample1() {

}
