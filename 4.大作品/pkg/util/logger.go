package util

import (
	"github.com/sirupsen/logrus"
	"os"
)

const logLevel = logrus.DebugLevel

var Logger *logrus.Logger


func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	Logger.SetLevel(logLevel)
}
