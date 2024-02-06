package shell

import (
	"context"

	"khranity/app/log"
	"khranity/app/lore"
	"khranity/app/utils"
	"khranity/app/shell/nix"
	"khranity/app/shell/win"
)

type Shell interface {
	Start(context.Context, *lore.Job) error
	Exec(context.Context, *lore.Job) error
	Get(context.Context, *lore.Job, string) error
	Put(context.Context, *lore.Job, string) error
}

// OS ...
type OS struct {
	Shell		Shell
}

var os OS

func New(ctx context.Context, logger *log.Logger) (*OS, error){
	switch(lore.Get().Setup.OS) {
	case "win", "windows":
		os.Shell 	= win.NewShell(logger)
	case "linux", "nix", "unix":
		os.Shell 	= nix.NewShell(logger)
	default:
		return nil, utils.ErrUndefinedOS
	}
	
	return &os, nil
}

func NewOS(ctx context.Context, logger *log.Logger, shell string) (*OS, error){
	switch(shell) {
	case "win", "windows":
		os.Shell 	= win.NewShell(logger)
	case "linux", "nix", "unix":
		os.Shell 	= nix.NewShell(logger)
	default:
		return nil, utils.ErrUndefinedOS
	}
	
	return &os, nil
}