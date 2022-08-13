package utils

import "github.com/threadedstream/wasmexperiments/internal/pkg/reporter"

func Assert(cond bool, msg string) {
	if !cond {
		reporter.ReportError(msg)
	}
}
