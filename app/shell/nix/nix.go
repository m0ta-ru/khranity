package nix

import (
	"khranity/app/log"
)

type ShellNix struct {
	logger 	*log.Logger
}

// NewShell ...
func NewShell(logger *log.Logger) *ShellNix {
	return &ShellNix{
		logger: logger,
	}
}