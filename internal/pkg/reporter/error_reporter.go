package reporter

import (
	"fmt"
	"os"
)

func ReportError(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}
