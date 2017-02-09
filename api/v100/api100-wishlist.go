package api100

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/carbocation/interpose"
	"github.com/chaosvermittlung/funkloch-server/db/v100"
	"github.com/gorilla/mux"
)

func getWishlistRouter(prefix string) *interpose.Middleware {
	r, m := GetNewSubrouter(prefix)
	r.HandleFunc("/", listWishlistsHandler).Methods("GET")
	r.HandleFunc("/", postWishlistHandler).Methods("POST")
	r.HandleFunc("/{ID}", getWishlistHandler).Methods("GET")
	r.HandleFunc("/{ID}", patchWishlistHandler).Methods("PATCH")
	r.HandleFunc("/{ID}", deleteWishlistHandler).Methods("DELETE")
	r.HandleFunc("/{ID}/Items", getWishlistItemsHandler).Methods("GET")
	r.HandleFunc("/{ID}/Item/{IID}/{Count}", addWishlistItemHandler).Methods("POST")
	r.HandleFunc("/{ID}/Item/{IID}", removeWishlistItemHandler).Methods("DELETE")
	return m
}

func postWishlistHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var wi db100.Wishlist
	err = decoder.Decode(&wi)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	err = wi.Insert()
	if err != nil {
		apierror(w, r, "Error Inserting Wishlist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&wi)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func listWishlistsHandler(w http.ResponseWriter, r *http.Request) {
	ww, err := db100.GetWishlists()
	if err != nil {
		apierror(w, r, "Error fetching Wishinglists: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&ww)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func getWishlistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	wi := db100.Wishlist{WishlistID: id}
	err = wi.GetDetails()
	if err != nil {
		apierror(w, r, "Error fetching Wishlist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&wi)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func patchWishlistHandler(w http.ResponseWriter, r *http.Request) {
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
	var wi db100.Wishlist
	err = decoder.Decode(&wi)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, ERROR_JSONERROR)
		return
	}
	wi.WishlistID = id
	err = wi.Update()
	if err != nil {
		apierror(w, r, "Error updating Wishlist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	j, err := json.Marshal(&wi)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func deleteWishlistHandler(w http.ResponseWriter, r *http.Request) {
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
	wi := db100.Wishlist{WishlistID: id}
	err = wi.Delete()
	if err != nil {
		apierror(w, r, "Error deleting Wishlist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
}

func addWishlistItemHandler(w http.ResponseWriter, r *http.Request) {
	err := userhasrRight(r, db100.USERRIGHT_MEMBER)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusUnauthorized, ERROR_USERNOTAUTHORIZED)
		return
	}
	vars := mux.Vars(r)
	i := vars["ID"]
	ii := vars["IID"]
	c := vars["Count"]
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
	count, err := strconv.Atoi(c)
	if err != nil {
		apierror(w, r, "Error converting Count: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	wli := db100.Wishlistitem{WishlistID: id, EquipmentID: iid, Count: count}

	err = wli.Insert()
	if err != nil {
		apierror(w, r, "Error inserting Item to Wishlist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}

}

func removeWishlistItemHandler(w http.ResponseWriter, r *http.Request) {
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
	wli := db100.Wishlistitem{WishlistID: id, EquipmentID: iid}

	err = wli.Delete()
	if err != nil {
		apierror(w, r, "Error deleting Item from Wishlist: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}

}

func getWishlistItemsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i := vars["ID"]
	id, err := strconv.Atoi(i)
	if err != nil {
		apierror(w, r, "Error converting ID: "+err.Error(), http.StatusBadRequest, ERROR_INVALIDPARAMETER)
		return
	}
	wi := db100.Wishlist{WishlistID: id}
	ee, err := wi.GetItems()
	if err != nil {
		apierror(w, r, "Error fetching Wishlistitems: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
		return
	}
	var wir []wishlistItemsResponse
	for _, e := range ee {
		wli := db100.Wishlistitem{WishlistID: wi.WishlistID, EquipmentID: e.EquipmentID}
		err := wli.GetDetails()
		if err != nil {
			apierror(w, r, "Error fetching Wishlistitem Details: "+err.Error(), http.StatusInternalServerError, ERROR_DBQUERYFAILED)
			return
		}
		w := wishlistItemsResponse{Equipment: e, Count: wli.Count}
		wir = append(wir, w)
	}
	j, err := json.Marshal(&wir)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
