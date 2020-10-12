package main

import (
	"fmt"

	nats "github.com/nats-io/nats.go"
)

//Ala is cool
func Ala() {
	fmt.Println("mod5")
	fmt.
		Println(nats.Version)
}
func main() {
	Ala()
}
