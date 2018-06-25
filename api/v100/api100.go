package api100

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Chaosvermittlung/funkloch-server/db/v100"
	"github.com/Chaosvermittlung/funkloch-server/global"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	jwt "gopkg.in/dgrijalva/jwt-go.v2"
)

const mySigningKey = "K3gQXQ4Xp87jERnQqYX3q6vyQZDrPZBYEDXVp6aPm78VD3S7wuxD2LB4VKX8S58sCEFwdybD"

func GetSubrouter(prefix string) *interpose.Middleware {
	middle100 := interpose.New()
	//middle800.Use(apiglobal.LoggerMiddleware())

	a100 := mux.NewRouter().PathPrefix(prefix).Subrouter()
	a100 = a100.StrictSlash(true)
	a100.HandleFunc("/auth", authHandler).Methods("GET")
	a100.HandleFunc("/auth-refesh", authRefreshHandler).Methods("GET")
	a100user := getUserRouter(prefix + "/user")
	a100.PathPrefix("/user").Handler(a100user)

	a100item := getItemRouter(prefix + "/item")
	a100.PathPrefix("/item").Handler(a100item)

	a100box := getBoxRouter(prefix + "/box")
	a100.PathPrefix("/box").Handler(a100box)

	a100store := getStoreRouter(prefix + "/store")
	a100.PathPrefix("/store").Handler(a100store)

	a100equip := getEquipmentRouter(prefix + "/equipment")
	a100.PathPrefix("/equipment").Handler(a100equip)

	a100event := getEventRouter(prefix + "/event")
	a100.PathPrefix("/event").Handler(a100event)

	a100fault := getFaultRouter(prefix + "/fault")
	a100.PathPrefix("/fault").Handler(a100fault)

	a100packinglist := getPackinglistRouter(prefix + "/packinglist")
	a100.PathPrefix("/packinglist").Handler(a100packinglist)

	a100wishlist := getWishlistRouter(prefix + "/wishlist")
	a100.PathPrefix("/wishlist").Handler(a100wishlist)

	middle100.UseHandler(a100)
	return middle100
}

func notfoundHandler(w http.ResponseWriter, r *http.Request) {
	apierror(w, r, "No Handler found for: "+r.RequestURI, http.StatusNotFound, ERROR_NOTFOUND)
}

func generateNewToken(un db100.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["iss"] = "funkloch"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token.Claims["user"] = un.UserID
	token.Claims["rights"] = un.Right
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString([]byte(tokenString)), err
}

func getTokenfromRequest(r *http.Request) (*jwt.Token, error) {
	m, t, err := getAuthorization(r)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(m) != "token" {
		return nil, errors.New("Token head missing")
	}

	data, err := base64.StdEncoding.DecodeString(t)
	if err != nil {
		return nil, err
	}
	fmt.Println("data", string(data))

	token, err := jwt.Parse(string(data), func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(mySigningKey), nil

	})

	return token, err
}

func getUserfromToken(token *jwt.Token) (db100.User, error) {
	un := db100.User{}

	ui, ok := token.Claims["user"].(float64)
	if !ok {
		return un, errors.New("No id")
	}

	uid := int(ui)
	un.UserID = uid
	un.GetDetails()

	return un, nil
}

func GetNewSubrouter(prefix string) (*mux.Router, *interpose.Middleware) {
	m := interpose.New()
	//m.Use(apiglobal.LoggerMiddleware())
	//	m.Use(authMiddleware())

	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	r = r.StrictSlash(true)
	r.NotFoundHandler = http.HandlerFunc(notfoundHandler)
	m.UseHandler(r)

	return r, m
}

func apierror(w http.ResponseWriter, r *http.Request, err string, httpcode int, ecode APIErrorcode) {
	//Erzeugt einen json error Response und gibt ihn über http.Error zurück
	log.Println(err)
	er := ErrorResponse{strconv.Itoa(httpcode), strconv.Itoa(int(ecode)), ecode.String() + ":" + err}
	j, erro := json.Marshal(&er)
	if erro != nil {
		return
	}
	//apiglobal.Apilog(err, swsglobal.LOGLEVEL_ERROR)
	http.Error(w, string(j), httpcode)
}

func getAuthorization(r *http.Request) (string, string, error) {
	auth := r.Header.Get("Authorization")
	s := strings.Split(auth, " ")
	if len(s) < 2 {
		return "", "", errors.New("Authorization header malformed. Expected \"Authorization token\" got " + auth)
	}
	return s[0], s[1], nil
}

func authMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := getTokenfromRequest(r)

			if err == nil && token.Valid {
				next.ServeHTTP(w, r)
			} else {
				m := ""
				if err != nil {
					m = err.Error()
					log.Println(err)
				} else {
					m = "Token invalid"
					log.Println("Tolen invalid")
				}
				apierror(w, r, m, 401, ERROR_MALFORMEDAUTH)
			}
		})
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {

	m, t, err := getAuthorization(r)
	if err != nil {
		apierror(w, r, err.Error(), 401, ERROR_MALFORMEDAUTH)
		return
	}

	if strings.ToLower(m) != "basic" {
		apierror(w, r, "Auth Request malformed", 401, ERROR_MALFORMEDAUTH)
		return
	}

	data, err := base64.StdEncoding.DecodeString(t)
	if err != nil {
		apierror(w, r, err.Error(), 401, ERROR_MALFORMEDAUTH)
		return
	}

	s := strings.Split(string(data), ":")
	if len(s) < 2 {
		apierror(w, r, "Auth Request malformed", 401, ERROR_MALFORMEDAUTH)
		return
	}

	u, p := s[0], s[1]

	b, err := db100.DoesUserExist(u)
	if err != nil {
		apierror(w, r, err.Error(), 500, ERROR_DBQUERYFAILED)
		return
	}
	if !b {
		apierror(w, r, "Wrong Username or Password", 401, ERROR_WRONGCREDENTIALS)
		return
	}

	un := db100.User{Username: u}
	err = un.GetDetailstoUsername()
	if err != nil {
		apierror(w, r, err.Error(), 500, ERROR_DBQUERYFAILED)
		return
	}

	pw, err := global.GeneratePasswordHash(p, un.Salt)
	if err != nil {
		apierror(w, r, err.Error(), 500, ERROR_NOHASH)
		return
	}
	if pw != un.Password {
		apierror(w, r, "Wrong Username or Password", 401, ERROR_WRONGCREDENTIALS)
		return
	}

	tokenString, err := generateNewToken(un)
	if err != nil {
		apierror(w, r, err.Error(), 500, ERROR_NOTOKEN)
		return
	}
	fmt.Println(tokenString)
	ar := authResponse{tokenString}
	j, err := json.Marshal(&ar)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func authRefreshHandler(w http.ResponseWriter, r *http.Request) {
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

	tokenString, err := generateNewToken(un)
	if err != nil {
		apierror(w, r, err.Error(), 500, ERROR_NOTOKEN)
		return
	}

	ar := authResponse{tokenString}
	j, err := json.Marshal(&ar)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, ERROR_JSONERROR)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func userhasrRight(r *http.Request, ri db100.UserRight) error {
	token, err := getTokenfromRequest(r)
	if err != nil {
		return err
	}

	ou, err := getUserfromToken(token)
	if err != nil {
		return err
	}
	if ou.Right < ri {
		return errors.New("User not permitted for this Action")
	}

	return nil
}
