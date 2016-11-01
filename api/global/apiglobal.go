package apiglobal

import (
	"fmt"
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
)

func GetSubrouter(prefix string) *mux.Router {
	globalmiddle := interpose.New()

	ag := mux.NewRouter().PathPrefix(prefix).Subrouter()
	ag.HandleFunc("/", apiHandler)

	globalmiddle.UseHandler(ag)
	return ag
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "To be done")
}
