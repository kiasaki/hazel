package data

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Teams    []Team `json:"teams,omitempty"`
}
