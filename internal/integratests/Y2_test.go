package integratests

import (
	"emobile/internal/models"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func (suite *TS) Test_02() {
	tests := []struct {
		name   string
		hand   string
		sub    models.Subscription
		status int
		noErr  bool
	}{
		{
			name:   "Drop",
			hand:   "/delete",
			sub:    models.Subscription{},
			noErr:  true,
			status: http.StatusOK,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			var resp *resty.Response
			var err error
			httpc := resty.New().SetBaseURL(suite.host)
			req := httpc.R().SetHeader("Content-Type", "application/json").SetDoNotParseResponse(false).
				SetBody(tt.sub)

			switch tt.hand {
			case "/delete":
				resp, err = req.Delete("/delete")
			}

			if tt.noErr {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
			suite.Require().Equal(tt.status, resp.StatusCode())



		})
	}
}
