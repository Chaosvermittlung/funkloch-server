package api100

import (
	"encoding/json"
	"net/http"
	"strconv"

	db100 "github.com/Chaosvermittlung/funkloch-server/pkg/db/v100"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
)

func getBoxRouter(prefix string) *interpose.Middleware {
	r, m := GetNewSubrouter(prefix)
	r.HandleFunc("/", postBoxHandler).Methods("POST")
	r.HandleFunc("/list", listBoxesHandler).Methods("GET")
	r.HandleFunc("/{ID}", getBoxHandler).Methods("GET")
	r.HandleFunc("/{ID}", patchBoxHandler).Methods("PATCH")
	r.HandleFunc("/{ID}", deleteBoxHandler).Methods("DELETE")
	r.HandleFunc("/{ID}/items", getBoxItemsHandler).Methods("GET")
	r.HandleFunc("/{ID}/items/{IID}", addItemtoBoxHandler).Methods("POST")
	r.HandleFunc("/{ID}/items/{IID}", removeItemfromBoxHandler).Methods("DELETE")
	return m
}

func convertBoxListEntryinBoxResponse(b db100.BoxlistEntry) boxResponse {
	var br boxResponse
	br.Box.BoxID = b.BoxID
	br.Box.Code = b.Code
	br.Box.Description = b.Description
	br.Box.Weight = b.Weight
	br.Store.StoreID = b.StoreID
	br.Store.Adress = b.Adress
	br.Store.ManagerID = b.ManagerID
	br.Store.Name = b.Name
	br.User.UserID = b.ManagerID
	br.User.Username = b.Username
	br.User.Email = b.Email
	br.User.Right = db100.UserRight(b.Right)
	return br
}

func convertBoxListinBoxResponseList(bb []db100.BoxlistEntry) []boxResponse {
	var res []boxResponse
	for _, b := range bb {
		br := convertBoxListEntryinBoxResponse(b)
		res = append(res, br)
	}
	return res
}

func postBoxHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var b db100.Box
	err = decoder.Decode(&b)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	err = b.Insert()
	if err != nil {
		apierror(w, r, "Error Inserting Box: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&b)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func listBoxesHandler(w http.ResponseWriter, r *http.Request) {
	bb, err := db100.GetBoxesJoined()
	if err != nil {
		apierror(w, r, "Error fetching Boxes: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	res := convertBoxListinBoxResponseList(bb)
	j, err := json.Marshal(&res)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getBoxHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	var b db100.Box
	b.BoxID = id
	ble, err := b.GetFullDetails()
	if err != nil {
		apierror(w, r, "Error fetching Box: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	br := convertBoxListEntryinBoxResponse(ble)
	j, err := json.Marshal(&br)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func patchBoxHandler(w http.ResponseWriter, r *http.Request) {
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
	var b db100.Box
	err = decoder.Decode(&b)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	b.BoxID = id
	err = b.Update()
	if err != nil {
		apierror(w, r, "Error updating Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&b)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func deleteBoxHandler(w http.ResponseWriter, r *http.Request) {
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
	b := db100.Box{BoxID: id}
	err = b.Delete()
	if err != nil {
		apierror(w, r, "Error deleting StoreItem: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
}

func getBoxItemsHandler(w http.ResponseWriter, r *http.Request) {
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
	var b db100.Box
	b.BoxID = id

	ile, err := b.GetBoxItemsJoined()
	if err != nil {
		apierror(w, r, "Error getting Box Items: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	ir := convertItemListinItemResponseList(ile)
	j, err := json.Marshal(&ir)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func addItemtoBoxHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting Box ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	ii := vars["IID"]
	iid, err := strconv.Atoi(ii)
	if err != nil {
		apierror(w, r, "Error converting Item ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	it := db100.Item{ItemID: iid}
	err = it.SetBox(id)
	if err != nil {
		apierror(w, r, "Error updating Item: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
}

func removeItemfromBoxHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	vars := mux.Vars(r)
	ii := vars["IID"]
	iid, err := strconv.Atoi(ii)
	if err != nil {
		apierror(w, r, "Error converting Item ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	it := db100.Item{ItemID: iid}
	err = it.SetBox(0)
	if err != nil {
		apierror(w, r, "Error updating Item: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
}
