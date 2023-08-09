package main_test

import (
	"log/syslog"

	myHook "testdata/mod5/hooks"

	slhooks "github.com/sirupsen/logrus/hooks/syslog"

	"github.com/kr/pretty" // this import is shared in test files and non-test files
	"github.com/spf13/jwalterweatherman"
)

func Example_HookAPi() {
	_ = jwalterweatherman.TRACE
	_, _ = slhooks.NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")

	pretty.Sprint("wow")

	myHook.HookAPi()
	// Output:
	// You have attempted to use a feature that is not yet implemented.
}
