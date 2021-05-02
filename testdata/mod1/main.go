package main

import (
	"fmt"
	"testdata/mod1/api"
	"testdata/mod1/version"
)

func main() {
	fmt.Println("mod1")

	v := version.Ver()
	fmt.Println("testdata/mod1/version.Ver: ", v)

	apiMsg, err := api.Api()
	fmt.Printf("testdata/mod1/api.apiMsg: %v, testdata/mod1/api.error: %v", apiMsg, err)

}
