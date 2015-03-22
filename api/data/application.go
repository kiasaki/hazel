package data

type Application struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StackID   string `json:"stack_id"`
	OwnerID   string `json:"owner_id"`
	OwnerType string `json:"owner_type"`
	GitURL    string `json:"git_url"`
}

func (_ Application) Bucket() []byte {
	return []byte("applications")
}

func init() {
	registerModel(Application{})
}
