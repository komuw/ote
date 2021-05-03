## ote          

![ote ci](https://github.com/komuw/ote/workflows/ote%20ci/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/komuw/ote/branch/main/graph/badge.svg)](https://codecov.io/gh/komuw/ote)
[![PkgGoDev](https://pkg.go.dev/badge/https://pkg.go.dev/github.com/komuw/ote)](https://pkg.go.dev/github.com/komuw/ote)
[![Go Report Card](https://goreportcard.com/badge/github.com/komuw/ote)](https://goreportcard.com/report/github.com/komuw/ote)


`ote` updates a packages' `go.mod` file with a comment next to all dependencies that are test dependencies; identifying them as such.   

It's name is derived from Kenyan hip hop artiste, `Oteraw`(One third of the hiphop group `Kalamashaka`).                               

By default, `go` and its related tools(`go mod` etc) do not differentiate regular dependencies from test ones, when updating/writing the `go.mod` file.    
There are various reasons why this is so, see [go/issues/26955](https://github.com/golang/go/issues/26955) & [go/issues/26913](https://github.com/golang/go/issues/26913)      
Thus `ote` fills that missing gap.   
It is not perfect, but it seems to work. See [How it works](#how-it-works)



## Installation

```shell
go install github.com/komuw/ote@latest
```           


## Usage
```bash
ote --help
```
```bash
examples:

    ote .                 # update go.mod in the current directory
        
    ote -f /tmp/myPkg     # update go.mod in the /tmp/myPkg directory

    ote -r                # write to stdout instead of updating the go.mod in the current directory

    ote -f /tmp/myPkg -r  # write to stdout instead of updating go.mod file in the /tmp/myPkg directory.  
```

If your application has a `go.mod` file like the following;
```bash
module github.com/pkg/myPkg

require (
	github.com/Shopify/sarama v1.26.4
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.0
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/stretchr/testify v1.6.1 // priorComment
	golang.org/x/mod v0.3.0
)
```
running `ote` will update the `go.mod` to the following;
```bash
module github.com/pkg/myPkg

require (
	github.com/Shopify/sarama v1.26.4
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.5.0 // test
	github.com/nats-io/nats-server/v2 v2.1.7 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/stretchr/testify v1.6.1 // test; priorComment
	golang.org/x/mod v0.3.0
)
```
ie; assuming that `github.com/google/go-cmp` and `github.com/stretchr/testify` are test-only dependencies in your application.


## Features.
- Update `go.mod` file with a comment `// test` next to any dependencies that are only used in tests.
- Makes only the minimal of changes to `go.mod` files.
- Preserves any prior comments that were in existence.
- If a dependency was a test-only dependency and then it starts been used in other non-test contexts, `ote` will also recognise that and remove the `// test` comment.


## How it works  
1. read `go.mod` file.
2. get all the imports of all the files used by the package    
  here we consider all the known build tags(`darwin`, `openbsd`, `riscv64` etc)    
3. get all the modules of which all the imports belong to.    
4. find which of those are test-only modules.   
5. update `go.mod` with a comment(`// test`) next to all the test dependencies.


## Inspiration(Hat tip)
1. [`x/mod`](https://pkg.go.dev/golang.org/x/mod).
2. [`go/packages`](https://pkg.go.dev/golang.org/x/tools/go/packages).
