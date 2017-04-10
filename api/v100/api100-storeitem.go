package api100

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carbocation/interpose"
	"github.com/chaosvermittlung/funkloch-server/db/v100"
	"github.com/gorilla/mux"
)

func getStoreItemRouter(prefix string) *interpose.Middleware {
	r, m := GetNewSubrouter(prefix)
	r.HandleFunc("/", postStoreItemHandler).Methods("POST")
	r.HandleFunc("/list", listStoreItemsHandler).Methods("GET")
	r.HandleFunc("/{ID}", getStoreItemHandler).Methods("GET")
	r.HandleFunc("/{ID}", patchStoreItemHandler).Methods("PATCH")
	r.HandleFunc("/{ID}", deleteStoreItemHandler).Methods("DELETE")
	return m
}

func postStoreItemHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var s db100.StoreItem
	err = decoder.Decode(&s)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	err = s.Insert()
	if err != nil {
		apierror(w, r, "Error Inserting Storeitem: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
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

func listStoreItemsHandler(w http.ResponseWriter, r *http.Request) {
	ss, err := db100.GetStoreItems()
	if err != nil {
		apierror(w, r, "Error fetching StoreItems: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	var res []storeItemResponse
	for _, s := range ss {
		var sir storeItemResponse
		sir.StoreItem = s
		sir.Equipment.EquipmentID = s.EquipmentID
		err = sir.Equipment.GetDetails()
		if err != nil {
			apierror(w, r, "Error fetching StoreItem Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		sir.Store.StoreID = s.StoreID
		err = sir.Store.GetDetails()
		if err != nil {
			apierror(w, r, "Error fetching StoreItem Store: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		res = append(res, sir)
	}
	j, err := json.Marshal(&res)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getStoreItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	var sir storeItemResponse
	sir.StoreItem.StoreItemID = id
	err = sir.StoreItem.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching StoreItem: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	sir.Store.StoreID = sir.StoreItem.StoreID
	err = sir.Store.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching StoreItem Store: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	sir.Equipment.EquipmentID = sir.StoreItem.EquipmentID
	err = sir.Equipment.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching StoreItem Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&sir)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func patchStoreItemHandler(w http.ResponseWriter, r *http.Request) {
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
	var si db100.StoreItem
	err = decoder.Decode(&si)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	si.StoreItemID = id
	err = si.Update()
	if err != nil {
		apierror(w, r, "Error updating Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&si)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func deleteStoreItemHandler(w http.ResponseWriter, r *http.Request) {
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
	s := db100.StoreItem{StoreItemID: id}
	err = s.Delete()
	if err != nil {
		apierror(w, r, "Error deleting StoreItem: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
}
