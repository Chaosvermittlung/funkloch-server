package api100

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carbocation/interpose"
	"github.com/Chaosvermittlung/funkloch-server/db/v100"
	"github.com/gorilla/mux"
)

func getPackinglistRouter(prefix string) *interpose.Middleware {
	r, m := GetNewSubrouter(prefix)
	r.HandleFunc("/", postPackinglistHandler).Methods("POST")
	r.HandleFunc("/list", listPackinglistsHandler).Methods("GET")
	r.HandleFunc("/{ID}", getPackinglistHandler).Methods("GET")
	r.HandleFunc("/{ID}", patchPackinglistHandler).Methods("PATCH")
	r.HandleFunc("/{ID}", deletePackinglistHandler).Methods("DELETE")
	r.HandleFunc("/{ID}/Items", getPackinglistItems).Methods("GET")
	r.HandleFunc("/{ID}/Item/{IID}", addPackinglistItemHandler).Methods("POST")
	r.HandleFunc("/{ID}/Item/{IID}", removePackinglistItemHandler).Methods("DELETE")
	return m
}

func postPackinglistHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var p db100.Packinglist
	err = decoder.Decode(&p)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	err = p.Insert()
	if err != nil {
		apierror(w, r, "Error Inserting Packinglist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&p)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func listPackinglistsHandler(w http.ResponseWriter, r *http.Request) {
	pp, err := db100.GetPackinglists()
	if err != nil {
		apierror(w, r, "Error fetching Packinglists: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&pp)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getPackinglistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	p := db100.Packinglist{PackinglistID: id}
	err = p.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Packinglist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&p)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func patchPackinglistHandler(w http.ResponseWriter, r *http.Request) {
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
	var pl db100.Packinglist
	err = decoder.Decode(&pl)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	pl.PackinglistID = id
	err = pl.Update()
	if err != nil {
		apierror(w, r, "Error updating Packinglist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&pl)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func deletePackinglistHandler(w http.ResponseWriter, r *http.Request) {
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
	p := db100.Packinglist{PackinglistID: id}
	err = p.Delete()
	if err != nil {
		apierror(w, r, "Error deleting Packinglist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
}

func addPackinglistItemHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	vars := mux.Vars(r)
	i := vars["ID"]
	ii := vars["IID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	iid, err := strconv.Atoi(ii)
	if err != nil {
		apierror(w, r, "Error converting IID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	pli := db100.PackinglistItem{PackinglistID: id, StoreitemID: iid}

	err = pli.Insert()
	if err != nil {
		apierror(w, r, "Error inserting Item to Packinglist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}

}

func removePackinglistItemHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	vars := mux.Vars(r)
	i := vars["ID"]
	ii := vars["IID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	iid, err := strconv.Atoi(ii)
	if err != nil {
		apierror(w, r, "Error converting IID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	pli := db100.PackinglistItem{PackinglistID: id, StoreitemID: iid}

	err = pli.Delete()
	if err != nil {
		apierror(w, r, "Error deleting Item from Packinglist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}

}

func getPackinglistItems(w http.ResponseWriter, r *http.Request) {
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
	p := db100.Packinglist{PackinglistID: id}
	sis, err := p.GetItems()
	if err != nil {
		apierror(w, r, "Error fetching Packinglist Items: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	var plir []packinglistItemsResponse
	for _, si := range sis {
		e := db100.Equipment{EquipmentID: si.EquipmentID}
		err := e.GetDetails()
		if err != nil {
			apierror(w, r, "Error fetching Equipment Details: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		s := db100.Store{StoreID: si.StoreID}
		err = s.GetDetails()
		if err != nil {
			apierror(w, r, "Error fetching Store Details: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		pli := packinglistItemsResponse{StoreItemID: si.StoreItemID, Equipment: e, Store: s}
		plir = append(plir, pli)
	}
	j, err := json.Marshal(&plir)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
