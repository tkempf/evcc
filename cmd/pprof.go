// +build prof

package cmd

import (
	"net/http"
	"net/http/pprof"
)

func init() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/pprof", pprof.Profile)
		log.ERROR.Println(http.ListenAndServe(":7077", mux))
	}()
}
