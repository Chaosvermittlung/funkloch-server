package api100

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/chaosvermittlung/funkloch-server/db/v100"
	"github.com/gorilla/mux"
)

func postStoreHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_ADMIN)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var s db100.Store
	err = decoder.Decode(&s)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	err = s.Insert()
	if err != nil {
		apierror(w, r, "Error Inserting Store: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&s)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getStoreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	s := db100.Store{ID: id}
	err = s.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Store: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&s)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getStoreManager(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	s := db100.Store{ID: id}
	u, err := s.GetManager()
	if err != nil {
		apierror(w, r, "Error fetching Store Manager: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&u)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func updateStoreHandler(w http.ResponseWriter, r *http.Request) {
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
	s := db100.Store{ID: id}
	decoder := json.NewDecoder(r.Body)
	var st db100.Store
	err = decoder.Decode(&st)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	err = s.Update()
	if err != nil {
		apierror(w, r, "Error updating Store: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&s)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func deleteStoreHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_ADMIN)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var s db100.Store
	err = decoder.Decode(&s)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	err = s.Delete()
	if err != nil {
		apierror(w, r, "Error deleting Store: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
}

func getStoreItemsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	s := db100.Store{ID: id}
	ii, err := s.GetStoreitems()
	if err != nil {
		apierror(w, r, "Error fetching Store Items: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}

	var res []db100.Equipment

	for _, si := range ii {
		e := db100.Equipment{ID: si.EquipmentID}
		err := e.GetDetails()
		if err != nil {
			apierror(w, r, "Error fetching Item Details: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		e.ID = si.ID
		res = append(res, e)
	}

	j, err := json.Marshal(&res)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
