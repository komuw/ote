package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func main() {
	fmt.Println("mod2/nestedModule1")

	errors.New("mod2/nestedModule1")
}
