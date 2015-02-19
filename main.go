package main

func main() {
	s := NewServer(ParseFlag())
	s.Run()
}
