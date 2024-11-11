package http

import (
	"net/http"
	"net/http/pprof"
)

func Register(hmux *http.ServeMux) {
	hmux.HandleFunc("/debug/pprof/", pprof.Index)
	hmux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	hmux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	hmux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	hmux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
