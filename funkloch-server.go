package main

import (
	"github.com/chaosvermittlung/funkloch-server/api/v100"

	"github.com/chaosvermittlung/funkloch-server/api/global"
	"github.com/chaosvermittlung/funkloch-server/db/v100"
	"github.com/chaosvermittlung/funkloch-server/global"
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
}
