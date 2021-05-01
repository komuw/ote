package main

import (
	"fmt"
	"testdata/mod2/api"
	"testdata/mod2/version"
)

func main() {
	fmt.Println("mod2")

	v := version.Ver()
	fmt.Println("testdata/mod2/version.Ver: ", v)

	apiMsg, err := api.Api()
	fmt.Printf("testdata/mod2/api.apiMsg: %v, testdata/mod2/api.error: %v", apiMsg, err)

}
