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
	"github.com/kiasaki/hazel/api/data"
)

var authHeaderPrefix = "Bearer "

type contextKey int

const AuthContextKey contextKey = 10

type LoginService struct {
	s *Server
}

func (e *LoginService) Register(router *mux.Router) {
	e.s.AddFilters(e.authMiddleware)

	e.s.Logger.Println("Registering [LoginService] at path [/auth]")
	router.HandleFunc("/auth/login", e.Login).Methods("POST")
}

// This is what is stored in the context and accessible from handlers
type Authorization struct {
	User  data.User
	Token jwt.Token
}

func (e *LoginService) authMiddleware(h http.Handler) http.Handler {
	keyFunc := func(_ *jwt.Token) (interface{}, error) {
		return []byte(e.s.Config.JwtSecret), nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer context.Clear(r)

		if !strings.HasPrefix(r.URL.Path, "/auth") {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, authHeaderPrefix) {
				rest.SetUnauthorizedResponse(w)
				return
			}

			tokenString := auth[len(authHeaderPrefix):]
			token, err := jwt.Parse(tokenString, keyFunc)

			if err == nil && token.Valid {
				var user = data.User{}
				err := e.s.DB.Get(token.Claims["u"].(string), &user)

				if err != nil {
					rest.SetUnauthorizedResponse(w)
					rest.WriteEntity(w, rest.J{"error": fmt.Sprint("Can't find user for this token:", err)})
					return
				} else {
					context.Set(r, AuthContextKey, Authorization{
						User:  user,
						Token: *token,
					})
				}
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

func TokenForUser(u data.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims["u"] = u.ID
	token.Claims["t"] = []string{"1"}
	token.Claims["name"] = u.FullName
	token.Claims["scope"] = "users applications builds"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(e.s.Config.JwtSecret))
	return tokenString, err
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (e *LoginService) Login(w http.ResponseWriter, r *http.Request) {
	var u = data.User{}
	tokenString, err := TokenForUser(u)
	if err != nil {
		rest.SetInternalServerErrorResponse(w, err)
		return
	}

	rest.SetOKResponse(w, LoginResponse{
		Token: tokenString,
	})
}

func (e *LoginService) Signup(w http.ResponseWriter, r *http.Request) {
	var user = data.User{}
	err := rest.Bind(r, &user)
	if err != nil || !user.Valid() {
		rest.SetBadRequestResponse(w)
		rest.WriteEntity(w, rest.J{"error": "A valid user must be passed with email, password and full_name"})
		return
	}

	err := user.HashPassword()
	if err != nil {
		rest.SetInternalServerError(w)
		return
	}
	err := e.s.DB.Save(user)
	if err != nil {
		rest.SetInternalServerError(w)
		return
	}
	token, err := TokenForUser(user)
	if err != nil {
		rest.SetInternalServerError(w)
		return
	}

	rest.SetOKResponse(w, SignupResponse{
		Token: token,
		User:  user,
	})
}
