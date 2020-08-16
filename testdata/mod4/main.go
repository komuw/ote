package main

import (
	"fmt"

	tomb1 "gopkg.in/tomb.v1"

	"gopkg.in/yaml.v2"

	tomb "gopkg.in/tomb.v2"
)

func main() {

	fmt.Println("mod4: ", tomb1.ErrStillAlive)

	fmt.Println("tomb mod4: ", tomb.Tomb{})

	_ = yaml.Encoder{}

}
