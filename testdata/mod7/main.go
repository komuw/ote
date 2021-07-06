package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus/collectors"
)

func main() {
	fmt.Println("mod7")
	_ = collectors.NewBuildInfoCollector()
}
