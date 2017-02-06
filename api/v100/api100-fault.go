package api100

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carbocation/interpose"
	"github.com/chaosvermittlung/funkloch-server/db/v100"
	"github.com/gorilla/mux"
)

func getFaultRouter(prefix string) *interpose.Middleware {
	r, m := GetNewSubrouter(prefix)
	r.HandleFunc("/", listFaultsHandler).Methods("GET")
	r.HandleFunc("/{ID}", getFaultHandler).Methods("GET")
	r.HandleFunc("/{ID}", patchFaultHandler).Methods("PATCH")
	r.HandleFunc("/{ID}", deleteFaultHandler).Methods("DELETE")
	return m
}

func listFaultsHandler(w http.ResponseWriter, r *http.Request) {
	ff, err := db100.GetFaults()
	if err != nil {
		apierror(w, r, "Error fetching Faults: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&ff)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getFaultHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	f := db100.Fault{FaultID: id}
	err = f.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Fault: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&f)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func patchFaultHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var fa db100.Fault
	err = decoder.Decode(&fa)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	fa.FaultID = id
	err = fa.Update()
	if err != nil {
		apierror(w, r, "Error updating Faul: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&fa)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func deleteFaultHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_ADMIN)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	s := db100.Fault{FaultID: id}
	err = s.Delete()
	if err != nil {
		apierror(w, r, "Error deleting Fault: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
}
