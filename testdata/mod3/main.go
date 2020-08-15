package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/go-cmp/cmp"

	"github.com/pkg/json"
)

func main() {
	a := 6
	b := 4 + 2
	if !cmp.Equal(a, b) {
		log.Fatalf("a is not equal to b")
	}

	fmt.Println("a,b: ", a, b)

	input := `{"a": 1,"b": 123.456, "c": [null]}`
	sc := json.NewScanner(strings.NewReader(input))

	fmt.Println("scanner: ", sc)
}
