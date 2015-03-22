package api

import (
	"net/http"
	"regexp"

	"github.com/kiasaki/batbelt/rest"
	"github.com/kiasaki/batbelt/uuid"
	"github.com/kiasaki/hazel/api/data"
)

type ApplicationsEndpoint struct {
	s *Server
}

func (e *ApplicationsEndpoint) Path() string {
	return "/applications"
}

func (e *ApplicationsEndpoint) GET(w http.ResponseWriter, r *http.Request) {
	rest.SetOKResponse(w, rest.J{"applications": []string{}})
}

type ApplicationCreateRequest struct {
	Name    string `json:"name"`
	GitURL  string `json:"git_url"`
	StackID string `json:"stack_id"`
}

func (e *ApplicationsEndpoint) POST(w http.ResponseWriter, r *http.Request) {
	var request ApplicationCreateRequest
	err := rest.Bind(r, &request)
	if err != nil || request.GitURL == "" || request.StackID == "" {
		rest.SetBadRequestResponse(w)
		rest.WriteEntity(w, rest.J{"error": "Application creation request requires name, git_url and stack_id"})
		return
	}

	// Verify name follows slug format
	match, err := regexp.Match("[a-z0-9-]{3,30}", []byte(request.Name))
	if err != nil || !match {
		rest.SetBadRequestResponse(w)
		rest.WriteEntity(w, rest.J{"error": "Application name must be 3 to 30 characters long and only composed from lower case characters, hyphens or numbers"})
		return
	}

	// Fetch current user, we need it for app owner
	auth, ok := GetRequestAuthorization(r)
	if !ok {
		rest.SetInternalServerErrorResponse(w, nil)
		return
	}

	// Create app struct and generate UUID for it
	app := data.Application{
		ID:        uuid.NewUUID().String(),
		Name:      request.Name,
		GitURL:    request.GitURL,
		StackID:   request.StackID,
		OwnerID:   auth.User.ID,
		OwnerType: "user",
	}

	// Save new app
	err = e.s.DB.Save(app.ID, app)
	if err != nil {
		rest.SetInternalServerErrorResponse(w, err)
		return
	}

	rest.SetOKResponse(w, app)
}
