package archiver

import (
	"context"

	"khranity/app/log"
	"khranity/app/lore"
	"khranity/app/archiver/mholt"
	"khranity/app/archiver/native"
	"khranity/app/archiver/targz"
	//"khranity/app/utils"
)

type Provider interface {
	Append(context.Context, *lore.Job, string, string) error
	Extract(context.Context, *lore.Job, string, string) error
}

// Archive ...
type Archive struct {
	Provider	Provider
}

var arch Archive

func New(ctx context.Context, logger *log.Logger, method string) (*Archive, error){
	switch(method) {
	case "native":
		arch.Provider 	= native.NewProvider(logger)
	case "targz", "lib-targz":
		arch.Provider 	= targz.NewProvider(logger)
	case "mholt", "lib-mholt":
		arch.Provider 	= mholt.NewProvider(logger)
	default:
		//return nil, utils.ErrUndefinedArchiver
		arch.Provider 	= native.NewProvider(logger)
	}
	
	return &arch, nil
}