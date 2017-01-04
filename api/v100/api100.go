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

	"github.com/carbocation/interpose"
	"github.com/chaosvermittlung/funkloch-server/db/v100"
	"github.com/gorilla/mux"
	jwt "gopkg.in/dgrijalva/jwt-go.v2"
)

const mySigningKey = "K3gQXQ4Xp87jERnQqYX3q6vyQZDrPZBYEDXVp6aPm78VD3S7wuxD2LB4VKX8S58sCEFwdybD"

func GetSubrouter(prefix string) *interpose.Middleware {
	middle100 := interpose.New()
	//middle800.Use(apiglobal.LoggerMiddleware())

	a100 := mux.NewRouter().PathPrefix(prefix).Subrouter()
	a100 = a100.StrictSlash(true)
	a100user := getUserRouter(prefix + "/user")
	a100.PathPrefix("/user").Handler(a100user)
	a100store := getStoreRouter(prefix + "/store")
	a100.PathPrefix("/store").Handler(a100store)

	middle100.UseHandler(a100)
	return middle100
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
	//	m.Use(apiglobal.LoggerMiddleware())
	//m.Use(authMiddleware())

	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	r = r.StrictSlash(true)
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

func userhasrRight(r *http.Request, ri db100.UserRight) error {
	token, _ := getTokenfromRequest(r)

	ou, err := getUserfromToken(token)
	if err != nil {
		return err
	}
	if ou.Right < ri {
		return errors.New("User not permitted for this Action")
	}

	return nil
}
