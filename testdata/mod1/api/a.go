package api

import (
	"github.com/hashicorp/nomad/drivers/shared/executor"
	"github.com/pkg/errors"

	"golang.org/x/sys/windows"

	//used in both test and non test files
	"rsc.io/quote"
)

func Api() (string, error) {
	quote.Glass()
	_ = executor.ExecutorVersionLatest
	_ = windows.EVENTLOG_SUCCESS

	return "hello from api", errors.New("not implemented")
}
