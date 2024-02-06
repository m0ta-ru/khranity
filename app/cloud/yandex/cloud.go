package yandex

import (
	"khranity/app/log"
	"khranity/app/lore"
)

type Cloud struct {
	logger 	*log.Logger
}

// NewCloud ...
func NewCloud(logger *log.Logger) *Cloud {
	return &Cloud{
		logger: logger,
	}
}

// TestCloud ...
func TestCloud(cloud *lore.Cloud) error {
	return nil
}
