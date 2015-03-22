package data

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (_ Team) Bucket() []byte {
	return []byte("teams")
}

func init() {
	registerModel(Team{})
}
