package statsviz

import (
	"github.com/arl/statsviz"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

func GoCoreMetrics(srv *http.Server) {
	r := mux.NewRouter()
	r.Methods("GET").Path("/debug/statsviz/ws").Name("GET /debug/statsviz/ws").HandlerFunc(statsviz.Ws)
	r.Methods("GET").PathPrefix("/debug/statsviz/").Name("GET /debug/statsviz/").Handler(statsviz.Index)
	srv.HandlePrefix("/", r)
}
