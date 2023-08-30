package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenStr, err := token.SignedString(SECRET)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return tokenStr, nil
}

func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("header: %s\n", r.Header)
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC) // This part of the code is using Go's type assertion feature to check if the method (algorithm) used for signing the JWT is HMAC (Hash-based Message Authentication Code).
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
				}
				return SECRET, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}
			if token.Valid {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}
	})
}

var SECRET = []byte("super-secret-auth-key")
var api_key = "123"

func GetJwt(w http.ResponseWriter, r *http.Request) {
	if r.Header["ACCESS"] != nil {
		if r.Header["ACCESS"][0] == api_key {
			token, err := CreateJWT()
			if err != nil {
				return
			}
			fmt.Fprint(w, token)
		}
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("super secret"))
}

func main() {
	http.Handle("/api", ValidateJWT(Home))
	http.HandleFunc("/jwt", GetJwt)

	http.ListenAndServe(":3500", nil)
}
