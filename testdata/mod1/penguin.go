// +build linux

package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func penguin() {
	fmt.Println(sarama.AclOperationAny)
}
