package win

import (
	"context"

	"khranity/app/lore"
	"khranity/app/log"
)

func (sh *Shell) Start(ctx context.Context, job *lore.Job) error {
	return nil
}

func (sh *Shell) Exec(ctx context.Context, job *lore.Job) error {
	sh.logger.Info("exec job", log.Object("job", job))

	return nil
}

func (sh *Shell) Get(ctx context.Context, job *lore.Job, temp string) error {
	sh.logger.Info("get job", log.Object("job", job))

	return nil
}

func (sh *Shell) Put(ctx context.Context, job *lore.Job, temp string) error {
	sh.logger.Info("put job", log.Object("job", job))

	return nil
}