package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Chaosvermittlung/funkloch-server/api/global"
	"github.com/Chaosvermittlung/funkloch-server/api/v100"
	"github.com/Chaosvermittlung/funkloch-server/db/v100"
	"github.com/Chaosvermittlung/funkloch-server/global"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	db100.Initialisation(&global.Conf.Connection)
	//API Handler
	//Setzt alle Routen zu den API Pfaden
	apig := apiglobal.GetSubrouter("/api")
	r.PathPrefix("/api").Handler(apig)
	//1.0.0 Api Version
	a100 := api100.GetSubrouter("/api/v100")
	apig.PathPrefix("/v100").Handler(a100)

	log.Println("funkloch Server Running")
	port := ":" + strconv.Itoa(global.Conf.Port)
	log.Fatal(http.ListenAndServe(port, r))
}
