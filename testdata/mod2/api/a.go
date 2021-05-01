package api

import (
	"github.com/hashicorp/nomad/drivers/shared/executor"
	"github.com/pkg/errors"

	//used in both test and non test files
	"rsc.io/quote"
)

func Api() (string, error) {
	quote.Glass()
	_ = executor.ExecutorVersionLatest
	return "hello from api", errors.New("not implemented")
}
