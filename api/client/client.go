package client

type Hazel struct {
	JwtToken string
}

func NewClient() Hazel {
	return Hazel{}
}
