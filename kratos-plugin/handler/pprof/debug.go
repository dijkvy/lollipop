package pprof

import (
	"net/http/pprof"
)
import "github.com/go-kratos/kratos/v2/transport/http"

func RegisterPprof(srv *http.Server) {
	srv.HandleFunc("/debug/pprof/", pprof.Index)
	srv.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	srv.HandleFunc("/debug/pprof/profile", pprof.Profile)
	srv.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	srv.HandleFunc("/debug/pprof/trace", pprof.Trace)
	srv.HandleFunc("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	srv.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	srv.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	srv.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	srv.HandleFunc("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
	srv.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}
