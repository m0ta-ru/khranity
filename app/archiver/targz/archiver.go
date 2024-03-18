package targz

import (
	"context"

	"khranity/app/lore"
)

func (pvd *Provider) Append(ctx context.Context, job *lore.Job, path string, file string) error {
	
	return Compress(path, file, job.Ignore)
}

func (pvd *Provider) Extract(ctx context.Context, job *lore.Job, file string, path string) error {

	return Extract(file, path)
}