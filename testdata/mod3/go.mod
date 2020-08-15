module testdata/mod3

go 1.14

// ote should remove the //test comment from go-cmp since it is also used in main.go
// it should also add a //test comment to testify
require (
	github.com/google/go-cmp v0.5.1 // test
	github.com/pkg/json v0.0.0-20200630040052-6ff993914616
	github.com/stretchr/testify v1.3.0
)
