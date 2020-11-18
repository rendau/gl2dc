package rest

import (
	"io/ioutil"
	"log"
	"net/http"
)

func (a *St) hMsg(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	a.cr.HandleMessage(bodyBytes)

	w.WriteHeader(200)
}
