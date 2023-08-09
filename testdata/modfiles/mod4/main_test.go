package main

import (
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/aws/aws-sdk-go/service/transfer"
	"github.com/go-kit/kit/log"
	zerolog "github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/errs/v2"
	"go.uber.org/zap"

	// this is the only test-only dependency
	"github.com/benweissmann/memongo"
)

func TestBaa(t *testing.T) {
	_ = transfer.DomainS3
	_ = errs.Tag("kk")
	_ = memongo.DBNameChars
	_ = scs.Session{}
	_ = logrus.FieldKeyMsg
	_ = log.DefaultTimestamp
	_ = zerolog.Logger
	_ = zap.DebugLevel
}
