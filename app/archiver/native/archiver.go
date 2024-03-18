package native

import (
	"context"
	
	"khranity/app/lore"
)

func (pvd *Provider) Append(ctx context.Context, job *lore.Job, path string, file string) error {
	err := tarAppend(path, file, job.Ignore)
	if (err != nil) {
		return err
	}
	return nil
}

func (pvd *Provider) Extract(ctx context.Context, job *lore.Job, file string, path string) error {
	err := tarExtract(file, path)
	if (err != nil) {
		return err
	}
	return nil
}