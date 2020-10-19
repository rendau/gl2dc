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

	err = a.cr.HandleMessage(bodyBytes)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(200)
}
