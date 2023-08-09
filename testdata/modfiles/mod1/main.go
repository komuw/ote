package main

import (
	"fmt"

	"testdata/modfiles/mod1/api"
	"testdata/modfiles/mod1/version"
)

/*
#include <stdlib.h>
*/
import "C"

func main() {
	fmt.Println("mod1")

	v := version.Ver()
	fmt.Println("testdata/modfiles/mod1/version.Ver: ", v)

	apiMsg, err := api.Api()
	fmt.Printf("testdata/modfiles/mod1/api.apiMsg: %v, testdata/modfiles/mod1/api.error: %v", apiMsg, err)

	_ = C.uint(45)
}
