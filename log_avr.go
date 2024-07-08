//go:build avr

package csp

import (
	"io"
)

var Logger io.Writer = nil

func log(_ string, args ...interface{}) {
	if Logger == nil {
		return
	}
	for _, arg := range args {
		if _, err := Logger.Write([]byte(arg.(string))); err != nil {
			return
		}
	}
}

func logTs(format string, args ...interface{}) {
	log(format, args...)
}
