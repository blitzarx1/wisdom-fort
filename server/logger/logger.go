package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	logFlags = log.Ldate | log.Ltime | log.Lmsgprefix

	prefixEnd       = ": "
	separatorLogger = "->"
)

func NewLogger(parentLogger *log.Logger, prefix string) *log.Logger {
	if parentLogger == nil {
		return log.New(os.Stdout, prefix+prefixEnd, logFlags)
	}

	trimmedParentPrefix := strings.TrimSuffix(parentLogger.Prefix(), prefixEnd)
	newPrefix := fmt.Sprintf("%s"+separatorLogger+"%s"+prefixEnd, trimmedParentPrefix, prefix)

	return log.New(parentLogger.Writer(), newPrefix, parentLogger.Flags())
}
