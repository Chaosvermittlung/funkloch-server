package api100

import (
	"encoding/json"
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/chaosvermittlung/funkloch-server/db/v100"
	"github.com/chaosvermittlung/funkloch-server/global"
	"github.com/gorilla/mux"
)

func getUserRouter(prefix string) *interpose.Middleware {
	r, m := GetNewSubrouter(prefix)
	r.HandleFunc("/", postUserHandler).Methods("POST")
	r.HandleFunc("/", patchCurrentUserHandler).Methods("PATCH")
	r.HandleFunc("/", getCurrentUserHandler).Methods("GET")
	r.HandleFunc("/list", listUsersHandler).Methods("GET")
	r.HandleFunc("/{name}", getUserHandler).Methods("GET")
	r.HandleFunc("/{name}", patchUserHandler).Methods("PATCH")
	r.HandleFunc("/{name}", deleteUserHandler).Methods("DELETE")

	return m
}

func getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenfromRequest(r)
	if err != nil {
		apierror(w, r, "Auth Request malformed", 401, ERROR_MALFORMEDAUTH)
		return
	}

	un, err := getUserfromToken(token)
	if err != nil {
		apierror(w, r, "Auth Request malformed", 401, ERROR_MALFORMEDAUTH)
		return
	}

	j, err := json.Marshal(&un)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func patchCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenfromRequest(r)
	if err != nil {
		apierror(w, r, "Auth Request malformed", 401, ERROR_MALFORMEDAUTH)
		return
	}

	ou, err := getUserfromToken(token)
	if err != nil {
		apierror(w, r, "Auth Request malformed", 401, ERROR_MALFORMEDAUTH)
	}

	contenttype := r.Header.Get("Content-Type")
	if contenttype != "application/json" {
		apierror(w, r, "Wrong contenttype. Expected: application/json Got: "+contenttype, http.StatusBadRequest, ERROR_FILEERROR)
		return
	}

	u := db100.User{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&u)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}

	if ou.UserID != u.UserID {
		apierror(w, r, "User not permitted for this Action", http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}

	ou.Patch(u)
	err = ou.Update()
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
	}

	j, err := json.Marshal(&ou)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func postUserHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_ADMIN)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var u db100.User
	err = decoder.Decode(&u)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}

	if u.Right < 1 {
		apierror(w, r, "User Right not set", http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}

	if u.Username == "" {
		apierror(w, r, "Username not set", http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}

	s, err := global.GenerateSalt()
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_NOHASH)
		return
	}
	u.Salt = s

	pw, err := global.GeneratePasswordHash(u.Password, u.Salt)
	if err != nil {
		apierror(w, r, err.Error(), 500, ERROR_NOHASH)
		return
	}
	u.Password = pw

	err = u.Insert()
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
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

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenfromRequest(r)
	if err != nil {
		apierror(w, r, "Auth Request malformed", 401, ERROR_MALFORMEDAUTH)
		return
	}

	ou, err := getUserfromToken(token)
	if err != nil {
		apierror(w, r, "Auth Request malformed", 401, ERROR_MALFORMEDAUTH)
		return
	}

	vars := mux.Vars(r)
	n := vars["name"]
	u := db100.User{Username: n}
	err = u.GetDetailstoUsername()
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}

	if (ou.UserID != u.UserID) && (ou.Right != db100.USERRIGHT_ADMIN) {
		apierror(w, r, "User not permitted for this Action", http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
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

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	uu, err := db100.GetUsers()
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	for i := range uu {
		uu[i].Password = ""
	}
	j, err := json.Marshal(&uu)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func patchUserHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_ADMIN)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}

	contenttype := r.Header.Get("Content-Type")
	if contenttype != "application/json" {
		apierror(w, r, "Wrong contenttype. Expected: application/json Got: "+contenttype, http.StatusBadRequest, ERROR_FILEERROR)
		return
	}

	vars := mux.Vars(r)
	n := vars["name"]
	ou := db100.User{Username: n}
	err = ou.GetDetailstoUsername()
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}

	u := db100.User{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&u)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}

	ou.Patch(u)
	err = ou.Update()
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}

	j, err := json.Marshal(&ou)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_ADMIN)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}

	vars := mux.Vars(r)
	n := vars["name"]
	u := db100.User{Username: n}
	err = u.GetDetailstoUsername()
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	err = db100.DeleteUser(u.UserID)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
}
