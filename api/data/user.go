package data

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Teams    []Team `json:"teams,omitempty"`
}

func (_ User) Bucket() []byte {
	return []byte("users")
}

func init() {
	registerModel(User{})
}

// Validates a user has an email, a full_name of some length and a password
func (u User) Valid() bool {
	return u.Email != "" && len(u.FullName) > 4 && len(u.Password) > 6
}
