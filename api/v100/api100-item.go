package api100

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Chaosvermittlung/funkloch-server/db/v100"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
)

func getItemRouter(prefix string) *interpose.Middleware {
	r, m := GetNewSubrouter(prefix)
	r.HandleFunc("/", postItemHandler).Methods("POST")
	r.HandleFunc("/list", listItemsHandler).Methods("GET")
	r.HandleFunc("/{ID}", getItemHandler).Methods("GET")
	r.HandleFunc("/{ID}", patchItemHandler).Methods("PATCH")
	r.HandleFunc("/{ID}", deleteItemHandler).Methods("DELETE")
	r.HandleFunc("/{ID}/fault", getItemFaultsHandler).Methods("GET")
	return m
}

func postItemHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var i db100.Item
	err = decoder.Decode(&i)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	err = i.Insert()
	if err != nil {
		apierror(w, r, "Error Inserting Item: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&i)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func listItemsHandler(w http.ResponseWriter, r *http.Request) {
	ss, err := db100.GetItemsJoined()
	if err != nil {
		apierror(w, r, "Error fetching StoreItems: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	var res []itemResponse
	for _, s := range ss {
		var sir itemResponse
		sir.Item.ItemID = s.ItemID
		sir.Item.Code = s.ItemCode
		sir.Item.EquipmentID = s.EquipmentID
		sir.Box.BoxID = s.BoxID
		sir.Box.Code = s.BoxCode
		sir.Box.Description = s.BoxDescription
		sir.Store.StoreID = s.StoreID
		sir.Store.Adress = s.StoreAddress
		sir.Store.ManagerID = s.StoreManagerID
		sir.Store.Name = s.StoreName
		sir.Equipment.EquipmentID = s.EquipmentID
		sir.Equipment.Name = s.EquipmentName
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

func getItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	var sir itemResponse
	sir.Item.ItemID = id
	err = sir.Item.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Item: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	sir.Box.BoxID = sir.Item.BoxID
	err = sir.Box.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Item Box: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	sir.Store.StoreID = sir.Box.StoreID
	err = sir.Store.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Item Box Store: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	sir.Equipment.EquipmentID = sir.Item.EquipmentID
	err = sir.Equipment.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Item Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
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

func patchItemHandler(w http.ResponseWriter, r *http.Request) {
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
	var si db100.Item
	err = decoder.Decode(&si)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	si.ItemID = id
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

func deleteItemHandler(w http.ResponseWriter, r *http.Request) {
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
	s := db100.Item{ItemID: id}
	err = s.Delete()
	if err != nil {
		apierror(w, r, "Error deleting StoreItem: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
}

func getItemFaultsHandler(w http.ResponseWriter, r *http.Request) {
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
	s := db100.Item{ItemID: id}
	ff, err := s.GetFaults()
	if err != nil {
		apierror(w, r, "Error getting Faults for Item: "+err.Error(), http.StatusBadRequest, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&ff)
	if err != nil {
		apierror(w, r, "Error Marshaling json: "+err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
