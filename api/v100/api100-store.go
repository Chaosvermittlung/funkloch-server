package api100

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/carbocation/interpose"
	"github.com/chaosvermittlung/funkloch-server/db/v100"
	"github.com/gorilla/mux"
)

func getStoreRouter(prefix string) *interpose.Middleware {
	r, m := GetNewSubrouter(prefix)
	r.HandleFunc("/", postStoreHandler).Methods("POST")
	r.HandleFunc("/list", listStoresHandler).Methods("GET")
	r.HandleFunc("/{ID}", getStoreHandler).Methods("GET")
	r.HandleFunc("/{ID}/Manager", getStoreManagerHandler).Methods("GET")
	r.HandleFunc("/{ID}", patchStoreHandler).Methods("PATCH")
	r.HandleFunc("/{ID}", deleteStoreHandler).Methods("DELETE")
	r.HandleFunc("/{ID}/Items", getStoreItemsHandler).Methods("GET")
	r.HandleFunc("/{ID}/ItemCount", getStoreItemCountHandler).Methods("GET")
	r.HandleFunc("/{ID}/NewItem", insertNewStoreItem).Methods("POST")

	return m
}

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

func listStoresHandler(w http.ResponseWriter, r *http.Request) {
	ss, err := db100.GetStores()
	if err != nil {
		apierror(w, r, "Error fetching Stores: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&ss)
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
	s := db100.Store{StoreID: id}
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

func getStoreManagerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	s := db100.Store{StoreID: id}
	err = s.GetDetails()
	if err != nil {
		apierror(w, r, "Error getting Store Detail: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
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

func patchStoreHandler(w http.ResponseWriter, r *http.Request) {
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
	decoder := json.NewDecoder(r.Body)
	var st db100.Store
	err = decoder.Decode(&st)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	st.StoreID = id
	err = st.Update()
	if err != nil {
		apierror(w, r, "Error updating Store: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&st)
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
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	s := db100.Store{StoreID: id}
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
	s := db100.Store{StoreID: id}
	ii, err := s.GetStoreitems()
	if err != nil {
		apierror(w, r, "Error fetching Store Items: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}

	var result []db100.Equipment

	for _, si := range ii {
		e := db100.Equipment{EquipmentID: si.EquipmentID}
		err := e.GetDetails()
		if err != nil {
			apierror(w, r, "Error fetching Item Details: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		e.EquipmentID = si.StoreItemID
		result = append(result, e)
	}

	j, err := json.Marshal(&result)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getStoreItemCountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	s := db100.Store{StoreID: id}
	ii, err := s.GetStoreitems()
	if err != nil {
		apierror(w, r, "Error fetching Store Items: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	m := make(map[int]int)

	for _, si := range ii {
		m[si.EquipmentID] = m[si.EquipmentID] + 1
	}

	var result []storeItemCountResponse
	for eid, ecount := range m {
		e := db100.Equipment{EquipmentID: eid}
		err := e.GetDetails()
		if err != nil {
			apierror(w, r, "Error fetching Item Details: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		stcr := storeItemCountResponse{Name: e.Name, Count: ecount}
		result = append(result, stcr)
	}

	j, err := json.Marshal(&result)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func insertNewStoreItem(w http.ResponseWriter, r *http.Request) {
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
	si.StoreID = id
	err = si.Insert()
	if err != nil {
		apierror(w, r, "Error while inserting Storeitem: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
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
