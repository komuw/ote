// +build darwin,cgo

package main

import (
	"fmt"

	nats "github.com/nats-io/nats.go"
)

func apple() {
	fmt.Println(nats.Version)
}
