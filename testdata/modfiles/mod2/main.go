package main

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

func main() {
	fmt.Println("mod2")
	_ = api.VerifyEchoRequest
}
