package sysinfo

import (
	"runtime"
)

func GoVersion() string {
	return runtime.Version()
}