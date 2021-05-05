package main_test

import (
	"log/syslog"
	myHook "testdata/mod5/hooks"

	slhooks "github.com/sirupsen/logrus/hooks/syslog"
)

func Example_HookAPi() {

	_, _ = slhooks.NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")

	myHook.HookAPi()
	// Output:
	// You have attempted to use a feature that is not yet implemented.
}
