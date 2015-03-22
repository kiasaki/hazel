package main

import "github.com/kiasaki/hazel/api"

func main() {
	s := api.NewServer()
	s.Setup()
	s.Run()
}
