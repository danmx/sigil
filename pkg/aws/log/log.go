package log

import (
	"github.com/aws/aws-sdk-go/aws"
	logr "github.com/sirupsen/logrus"
)

// A traceLogger provides a minimalistic logger satisfying the aws.Logger interface.
type traceLogger struct{}

// newTraceLogger returns a Logger which will write log messages to current logger
func NewTraceLogger() aws.Logger {
	return &traceLogger{}
}

// Log logs the parameters to the stdlib logger. See log.Println.
func (l traceLogger) Log(args ...interface{}) {
	logr.Trace(args...)
}
