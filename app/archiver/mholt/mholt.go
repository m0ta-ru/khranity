package mholt

import (
	"khranity/app/log"
)

type Provider struct {
	logger 	*log.Logger
}

// NewProvider ...
func NewProvider(logger *log.Logger) *Provider {
	return &Provider{
		logger: logger,
	}
}