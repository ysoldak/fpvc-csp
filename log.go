//go:build !avr

package csp

import (
	"fmt"
	"io"
	"time"
)

var Logger io.Writer = nil

func log(format string, args ...interface{}) {
	if Logger == nil {
		return
	}
	Logger.Write([]byte(fmt.Sprintf(format, args...)))
}

func logTs(format string, args ...interface{}) {
	if Logger == nil {
		return
	}
	args = append([]interface{}{time.Now().Format("15:04:05.000")}, args...)
	Logger.Write([]byte(fmt.Sprintf("%s "+format, args...)))
}
