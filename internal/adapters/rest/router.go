package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (a *St) createRouter() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("", a.hMsg)
	r.HandleFunc("/", a.hMsg)

	return http.Handler(r)
}
