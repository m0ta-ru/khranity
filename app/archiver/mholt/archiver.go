package mholt

import (
	"context"

	"khranity/app/log"
	"khranity/app/lore"
)

func (pvd *Provider) Append(ctx context.Context, job *lore.Job, path string, file string) error {
	pvd.logger.Debug("append archive", 
		log.Object("job", job),
		log.String("path", path),
		log.String("file", file),
	)

	return nil
}

func (pvd *Provider) Extract(ctx context.Context, job *lore.Job, file string, path string) error {
	pvd.logger.Debug("extract archive", 
		log.Object("job", job),
		log.String("file", file),
		log.String("path", path),
	)

	return nil
}