package api100

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carbocation/interpose"
	"github.com/chaosvermittlung/funkloch-server/db/v100"
	"github.com/gorilla/mux"
)

func getEquipmentRouter(prefix string) *interpose.Middleware {
	r, m := GetNewSubrouter(prefix)
	r.HandleFunc("/", postEquipmentHandler).Methods("POST")
	r.HandleFunc("/", listEquipmentHandler).Methods("GET")
	r.HandleFunc("/count", getEquipmentsCountHandler).Methods("GET")
	r.HandleFunc("/{ID}", getEquipmentHandler).Methods("GET")
	r.HandleFunc("/{ID}", deleteEquipmentHandler).Methods("DELETE")
	r.HandleFunc("/{ID}", patchEquipmentHandler).Methods("PATCH")
	r.HandleFunc("/{ID}/list", getEquipmentCountHandler).Methods("GET")

	return m
}

func postEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var e db100.Equipment
	err = decoder.Decode(&e)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	err = e.Insert()
	if err != nil {
		apierror(w, r, "Error Inserting Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&e)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func listEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	ee, err := db100.GetEquipment()
	if err != nil {
		apierror(w, r, "Error fetching Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&ee)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	e := db100.Equipment{EquipmentID: id}
	err = e.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&e)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func patchEquipmentHandler(w http.ResponseWriter, r *http.Request) {
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
	var equ db100.Equipment
	err = decoder.Decode(&equ)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	equ.EquipmentID = id
	err = equ.Update()
	if err != nil {
		apierror(w, r, "Error updating Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&equ)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func deleteEquipmentHandler(w http.ResponseWriter, r *http.Request) {
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
	e := db100.Equipment{EquipmentID: id}
	err = e.Delete()
	if err != nil {
		apierror(w, r, "Error deleting Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
}

func getEquipmentsCountHandler(w http.ResponseWriter, r *http.Request) {
	ee, err := db100.GetEquipment()
	if err != nil {
		apierror(w, r, "Error fetching Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	ss, err := db100.GetStores()
	if err != nil {
		apierror(w, r, "Error fetching Stores: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	var result []equipmentCountResponse

	for _, e := range ee {
		for _, s := range ss {
			count, err := s.GetItemCount(e.EquipmentID)
			if err != nil {
				apierror(w, r, "Error fetching ItemCount: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
				return
			}
			if count > 0 {
				var ecr equipmentCountResponse
				ecr.Equipment = e
				ecr.Store = s
				ecr.Count = count
				result = append(result, ecr)
			}
		}

	}

	j, err := json.Marshal(&result)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getEquipmentCountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	e := db100.Equipment{EquipmentID: id}
	err = e.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Equipment: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	ss, err := db100.GetStores()
	if err != nil {
		apierror(w, r, "Error fetching Stores: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	var result []equipmentCountResponse

	for _, s := range ss {
		count, err := s.GetItemCount(e.EquipmentID)
		if err != nil {
			apierror(w, r, "Error fetching ItemCount: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		if count > 0 {
			var ecr equipmentCountResponse
			ecr.Equipment = e
			ecr.Store = s
			ecr.Count = count
			result = append(result, ecr)
		}
	}

	j, err := json.Marshal(&result)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
