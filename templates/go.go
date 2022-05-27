package templates

const Go = `package main

func main() {
	socket, err := cg.NewSocket("%s")
	if err != nil {
		log.Fatalf("failed to connect to server: $s", err)
	}
}`
