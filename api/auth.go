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
	"github.com/kiasaki/batbelt/uuid"
	"github.com/kiasaki/hazel/api/data"
	"golang.org/x/crypto/bcrypt"
)

var authHeaderPrefix = "Bearer "

type contextKey int

const AuthContextKey contextKey = 10

type AuthService struct {
	s *Server
}

func (e *AuthService) Register(router *mux.Router) {
	e.s.AddFilters(e.authMiddleware)

	e.s.Logger.Println("Registering [AuthService] at path [/auth]")
	router.HandleFunc("/auth/login", e.Login).Methods("POST")
	router.HandleFunc("/auth/signup", e.Signup).Methods("POST")
}

// This is what is stored in the context and accessible from handlers
type Authorization struct {
	User  data.User
	Token jwt.Token
}

func GetRequestAuthorization(r *http.Request) Authorization {
	if a, ok := context.Get(r, AuthContextKey); ok {
		return a.(Authorization)
	}
	return nil
}

func (e *AuthService) authMiddleware(h http.Handler) http.Handler {
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

func TokenForUser(secret string, u data.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	var teamIds []string
	for _, team := range u.Teams {
		teamIds = append(teamIds, team.ID)
	}

	token.Claims["u"] = u.ID
	token.Claims["t"] = teamIds
	token.Claims["name"] = u.FullName
	token.Claims["scope"] = "users applications builds"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func (e *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	var request = LoginRequest{}
	err := rest.Bind(r, &request)
	if err != nil {
		rest.SetBadRequestResponse(w)
		rest.WriteEntity(w, rest.J{"error": "A login request needs an email and password to be passed in"})
		return
	}

	// Find the user associated with email address requested
	user, err := e.s.DB.GetUserByEmail(request.Email)
	if err == data.ErrEntityNotFound {
		rest.SetNotFoundResponse(w)
		rest.WriteEntity(w, rest.J{"error": "There is no user with that email address"})
		return
	} else if err != nil {
		rest.SetInternalServerErrorResponse(w, err)
		return
	}

	// Check users password against the one provided
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		rest.SetBadRequestResponse(w)
		rest.WriteEntity(w, rest.J{"error": "Wrong password or email address provided"})
		return
	}

	tokenString, err := TokenForUser(e.s.Config.JwtSecret, user)
	if err != nil {
		rest.SetInternalServerErrorResponse(w, err)
		return
	}

	rest.SetOKResponse(w, LoginResponse{
		Token: tokenString,
	})
}

type SignupResponse struct {
	Token string    `json:"token"`
	User  data.User `json:"user"`
}

func (e *AuthService) Signup(w http.ResponseWriter, r *http.Request) {
	var user = data.User{}
	err := rest.Bind(r, &user)
	if err != nil || !user.Valid() {
		rest.SetBadRequestResponse(w)
		rest.WriteEntity(w, rest.J{"error": "A valid user must be passed with email, password and full_name"})
		return
	}

	// Check for duplicate email address
	sameEmailUser, err := e.s.DB.GetUserByEmail(user.Email)
	if err != nil && err != data.ErrEntityNotFound {
		rest.SetInternalServerErrorResponse(w, err)
		return
	}
	// At this point either err is nil of we didn't find user, so check match
	if sameEmailUser.Email == user.Email {
		rest.SetBadRequestResponse(w)
		rest.WriteEntity(w, rest.J{"error": "Email address already taken"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 11)
	if err != nil {
		rest.SetInternalServerErrorResponse(w, err)
		return
	}
	user.ID = uuid.NewUUID().String()
	user.Password = string(hashedPassword)
	e.s.Logger.Println(user)
	err = e.s.DB.Save(user.ID, user)
	if err != nil {
		rest.SetInternalServerErrorResponse(w, err)
		return
	}
	token, err := TokenForUser(e.s.Config.JwtSecret, user)
	if err != nil {
		rest.SetInternalServerErrorResponse(w, err)
		return
	}

	// Filter out password
	user.Password = ""
	rest.SetOKResponse(w, SignupResponse{
		Token: token,
		User:  user,
	})
}
