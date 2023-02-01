/*
Package chalker is a logging interface for the commands->stdout and uses the chalk package
*/
package chalker

import (
	"errors"

	"github.com/fatih/color"
)

// Logging color/style types
const (
	DEFAULT = "default"
	DIM     = "dim"
	BOLD    = "bold"
	ERROR   = "error"
	INFO    = "info"
	SUCCESS = "success"
	WARN    = "warn"
)

// Error chalks and returns an error
func Error(body string) error {
	return errors.New(color.MagentaString(body))
}

// Log writes chalks to console
func Log(level string, body string) {
	switch level {
	case INFO:
		color.Cyan(body)
	case WARN:
		color.Yellow(body)
	case ERROR:
		color.Magenta(body)
	case SUCCESS:
		color.Green(body)
	case DEFAULT:
		fallthrough
	default:
		color.White(body)
	}
}
