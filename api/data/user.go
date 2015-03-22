package data

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
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

func (d *Database) GetUserByEmail(email string) (User, error) {
	var user User
	err := d.DB.View(func(tx *bolt.Tx) error {
		return tx.Bucket(user.Bucket()).ForEach(func(k, v []byte) error {
			var u User
			err := json.Unmarshal(v, &u)
			if err != nil {
				return err
			}
			if u.Email == email {
				user = u
				// Stop iterating now, faster than always going through
				// all k/vs, error catched down feww lines
				return ErrEntityFound
			}
			return nil
		})
	})
	if err == ErrEntityFound {
		return user, nil
	} else if err != nil {
		return user, err
	} else {
		// If err was nil then we went through all entities and didn't match
		// on the provided email address
		return user, ErrEntityNotFound
	}
}
