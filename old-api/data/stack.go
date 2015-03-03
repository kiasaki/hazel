package data

import "github.com/kiasaki/batbelt/uuid"

type Stack struct {
	Id      string
	Name    string
	BaseAmi string
	Region  string
	VmSize  int
}

func (s Stack) Validate() bool {
	if len(s.Name) > 0 {
		return errors.New("Stack Name is required")
	}

	return true
}
