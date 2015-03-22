package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/kiasaki/batbelt/rest"
)

var authHeaderPrefix = "Bearer "

type contextKey int

const AuthContextKey contextKey = 10

type LoginService struct {
	s *Server
}

func (e *LoginService) Register(router *mux.Router) {
	e.s.Filters.Append(e.authMiddleware)

	e.s.Logger.Println("Registering [LoginService] at path [/auth]")
	router.HandleFunc("/auth/login", e.Login).Methods("POST")
}

func (e *LoginService) authMiddleware(h http.Handler) http.Handler {
	keyFunc := func(_ *jwt.Token) (interface{}, error) {
		return e.s.Config.JwtSecret, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer context.Clear(r)

		if !strings.HasPrefix(r.URL.Path, "/auth") {
			auth := r.Header.Get("Authentication")
			if !strings.HasPrefix(auth, authHeaderPrefix) {
				rest.SetUnauthorizedResponse(w)
				return
			}

			tokenString := auth[len(authHeaderPrefix):]
			token, err := jwt.Parse(tokenString, keyFunc)

			if err == nil && token.Valid {
				context.Set(r, AuthContextKey, *token)
			} else if ve, ok := err.(*jwt.ValidationError); ok {
				rest.SetUnauthorizedResponse(w)
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					rest.WriteEntity(w, rest.J{"error": "Token wasn't passed or is malformed"})
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					rest.WriteEntity(w, rest.J{"error": "Token expired or not valid yet"})
				} else {
					rest.WriteEntity(w, rest.J{"error": fmt.Sprint("Couldn't handle this token:", err)})
				}
				return
			} else {
				rest.SetUnauthorizedResponse(w)
				rest.WriteEntity(w, rest.J{"error": fmt.Sprint("Couldn't handle this token:", err)})
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (e *LoginService) Login(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims["u"] = "1"
	token.Claims["t"] = []string{"1"}
	token.Claims["name"] = "John"
	token.Claims["scope"] = "users applications builds"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(e.s.Config.JwtSecret))
	if err != nil {
		rest.SetInternalServerErrorResponse(w, err)
		return
	}

	rest.SetOKResponse(w, LoginResponse{
		Token: tokenString,
	})
}
