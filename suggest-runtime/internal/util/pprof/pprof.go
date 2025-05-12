package pprof

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/labstack/gommon/log"
)

func registerProfile(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}

func HandlePprof() {
	// pprof
	go func() {
		pprofPort := 3366
		log.Info(fmt.Sprintf("start pprof server on %d", pprofPort))
		pprofMux := http.NewServeMux()
		registerProfile(pprofMux)
		s := http.Server{
			Addr:           fmt.Sprintf("%s:%d", "0.0.0.0", pprofPort),
			Handler:        pprofMux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   5 * time.Minute,
			MaxHeaderBytes: 1 << 11,
		}
		log.Warn(s.ListenAndServe())
	}()
}
