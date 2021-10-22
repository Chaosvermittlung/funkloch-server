package api100

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Chaosvermittlung/funkloch-server/internal/global"
	db100 "github.com/Chaosvermittlung/funkloch-server/pkg/db/v100"
	"github.com/carbocation/interpose"
	jwt "github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

//Holds claims for JWT authentication
type FunklochClaims struct {
	User   int             `json:"user"`
	Rights db100.UserRight `json:"rights"`
	jwt.StandardClaims
}

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
	mySigningKey := global.Conf.TokenKey
	// Set some claims
	claims := FunklochClaims{
		un.UserID,
		un.Right,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "funkloch",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString([]byte(tokenString)), err
}

func getTokenfromRequest(r *http.Request) (*FunklochClaims, error) {
	mySigningKey := global.Conf.TokenKey
	m, t, err := getAuthorization(r)
	if err != nil {
		return nil, err
	}

	if m != "Bearer" {
		return nil, errors.New("bearer head missing")
	}

	data, err := base64.StdEncoding.DecodeString(t)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(string(data), &FunklochClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if err != nil {
		log.Println("Error parsing token: ", err)
		return nil, err
	}

	claims, ok := token.Claims.(*FunklochClaims)

	if !ok {
		err := errors.New("wrong claim format")
		log.Println(err)
		return nil, err
	}

	if !token.Valid {
		err := errors.New("token invalid")
		log.Println(err)
		return nil, err
	}
	return claims, err
}

func getUserfromToken(claims *FunklochClaims) (db100.User, error) {
	un := db100.User{}
	un.UserID = claims.User
	un.GetDetails()

	return un, nil
}

func GetNewSubrouter(prefix string) (*mux.Router, *interpose.Middleware) {
	m := interpose.New()
	//m.Use(apiglobal.LoggerMiddleware())
	m.Use(authMiddleware())

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
		return "", "", errors.New("Authorization header malformed. Expected \"Bearer <token>\" got " + auth)
	}
	return s[0], s[1], nil
}

func authMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := getTokenfromRequest(r)

			if err == nil {
				next.ServeHTTP(w, r)
			} else {
				apierror(w, r, "", 401, ERROR_MALFORMEDAUTH)
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
	claims, err := getTokenfromRequest(r)
	if err != nil {
		apierror(w, r, "Auth Request malformed", 401, ERROR_MALFORMEDAUTH)
		return
	}

	un, err := getUserfromToken(claims)
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
	claims, err := getTokenfromRequest(r)
	if err != nil {
		return err
	}

	ou, err := getUserfromToken(claims)
	if err != nil {
		return err
	}
	if ou.Right < ri {
		return errors.New("User not permitted for this Action")
	}

	return nil
}
