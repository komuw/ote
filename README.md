## ote          

[![ci](https://github.com/komuw/ote/workflows/ote%20ci/badge.svg)](https://github.com/komuw/ote/actions)
[![codecov](https://codecov.io/gh/komuw/ote/branch/master/graph/badge.svg)](https://codecov.io/gh/komuw/ote)


`ote` updates a packages' `go.mod` file with a comment next to all dependencies that are test dependencies; identifying them as such.   

It's name is derived from Kenyan hip hop artiste, `Oteraw`(One third of the hiphop group `Kalamashaka`).                               

By default, `go` and its related tools(`go mod` etc) do not differentiate regular dependencies from test ones when updating/writing the `go.mod` file.    
There are various reasons why this is so, see [go/issues/26955](https://github.com/golang/go/issues/26955) [go/issues/26913](https://github.com/golang/go/issues/26913)      
Thus `ote` fills that missing gap.   
It is not perfect, but it seems to work when it works. See [How it works](#how-it-works)


Comprehensive documetion is available -> [Documentation](https://pkg.go.dev/github.com/komuw/ote)


## Installation

```shell
go get github.com/komuw/ote
```           


## Usage
```bash
Usage of ote:
	ote .
		update go.mod file with a comment next to all dependencies that are test dependencies.
	ote -r .
		(readonly) display how the updated go.mod file would look like, without actually updating the file.
```

`ote` takes a `go.mod` file like;
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
and turns it into a `go.mod` file like;
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

## How it works  
1. read `go.mod` file.
2. get all the imports of all the files used by the package    
  here we consider all the known build tags(`darwin`, `openbsd`, `riscv64` etc)    
3. get all the modules of which all the imports belong to.    
4. find all the modules declared in `go.mod` file.   
5. find the difference in the modules between step 3 & 4.     
   this difference represents the test dependencies.  
6. update `go.mod` with a comment(`// test`) next to all the test dependencies.


## Inspiration(Hat tip)
1. [`x/mod`](https://pkg.go.dev/golang.org/x/mod).
2. [`go/packages`](https://pkg.go.dev/golang.org/x/tools/go/packages).
