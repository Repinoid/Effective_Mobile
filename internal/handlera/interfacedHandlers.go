package handlera

import "net/http"

type Handlers interface {
	DBPinger(rwr http.ResponseWriter, req *http.Request)
	CreateSub(rwr http.ResponseWriter, req *http.Request)
	ReadSub(rwr http.ResponseWriter, req *http.Request)
	ListSub(rwr http.ResponseWriter, req *http.Request)
	UpdateSub(rwr http.ResponseWriter, req *http.Request)
	DeleteSub(rwr http.ResponseWriter, req *http.Request)
	SumSub(rwr http.ResponseWriter, req *http.Request)
}
