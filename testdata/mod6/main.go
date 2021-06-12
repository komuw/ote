package main

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

func main() {
	fmt.Println("mod6")
	_ = api.VerifyEchoRequest
	errors.New("mod2/nestedModule1")
}
