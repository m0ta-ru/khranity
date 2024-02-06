package win

import (
	"khranity/app/log"
)

type Shell struct {
	logger 	*log.Logger
}

// NewShell ...
func NewShell(logger *log.Logger) *Shell {
	return &Shell{
		logger: logger,
	}
}