package main

import (
	"fmt"

	"github.com/DataDog/zstd"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func main() {
	fmt.Println("mod7")
	_ = collectors.NewBuildInfoCollector()
	_ = zstd.BestCompression
}
