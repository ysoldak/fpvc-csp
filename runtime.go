//go:build avr

package csp

import _ "unsafe" // for go:linkname to work

//go:linkname runtime_nanotime runtime.nanotime
func runtime_nanotime() int64
