package main

import (
	"context"
	"os"

	"khranity/app/log"
	"khranity/app/lore"
	"khranity/app/shell"
)

type GetCommand struct {
	Name  string `short:"n" long:"name" description:"Name of job to get"`
	Path  string `short:"p" long:"path" description:"Path to output"`
	Token string `short:"t" long:"token" description:"Token file to decrypt"`
}

var getCommand GetCommand

func (cmd *GetCommand) Execute(args []string) error {
	fLore 	:= parser.Groups()[0].FindOptionByLongName("lore").Value().(string)
	logger 	:= log.GetDefault()
	lore.Load(fLore)

	job, err := lore.GetJob(cmd.Name)
	if err != nil {
		return err
	}

	if len(cmd.Path) > 0 {
		job.Path = cmd.Path
	}

	return get(job, logger)
}

func init() {
	parser.AddCommand("get",
		"Get a object from cloud",
		"The \"get\" command gets a <name> from cloud to <path>.",
		&getCommand)
}

func get(job *lore.Job, logger *log.Logger) error {
	ctx := context.Background()

	s, err := shell.NewOS(ctx, logger, "nix")
	if err != nil {
		return err
	}

	return s.Shell.Get(ctx, job, os.TempDir())
}
