package main

import (
	"fmt"

	"github.com/alexedwards/scs/v2"
	"github.com/aws/aws-sdk-go/service/transfer"
	"github.com/go-kit/kit/log"
	"github.com/ishidawataru/sctp"
	zerolog "github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/errs/v2"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("mod4")
	// This is an application that share a lot of
	// dependencies between test files and non-test files.
	// see: https://github.com/komuw/ote/issues/22

	_ = sctp.MSG_NOTIFICATION
	_ = transfer.DomainEfs
	_ = errs.Tag("kk")
	_ = scs.GobCodec{}
	_ = logrus.FieldKeyMsg
	_ = log.DefaultTimestamp
	_ = zerolog.Logger
	_ = zap.ErrorLevel
}
