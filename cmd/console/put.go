package main

import (
	"context"
	"os"

	"khranity/app/log"
	"khranity/app/lore"
	"khranity/app/shell"
)

type PutCommand struct {
	Name  string `short:"n" long:"name" description:"Name of job to put"`
	Path  string `short:"p" long:"path" description:"Path to input"`
	Cloud string `short:"c" long:"cloud" description:"Cloud specified in lore-file"`
	Token string `short:"t" long:"token" description:"Token file to encrypt"`
}

var putCommand PutCommand

func (cmd *PutCommand) Execute(args []string) error {
	fLore 	:= parser.Groups()[0].FindOptionByLongName("lore").Value().(string)
	logger	:= log.GetDefault()
	lore.Load(fLore)

	job, err := lore.GetJob(cmd.Name)
	if err != nil {
		job = &lore.Job{}
		job.Name = cmd.Name
		job.Token = cmd.Token
		job.Cloud = cmd.Cloud
	}

	if len(cmd.Path) > 0 {
		job.Path = cmd.Path
	}

	return put(job, logger)
}

func init() {
	parser.AddCommand("put",
		"Put a object to cloud",
		"The \"put\" command puts <path> to cloud by <name>.",
		&putCommand)
}

func put(job *lore.Job, logger *log.Logger) error {
	ctx := context.Background()

	s, err := shell.NewOS(ctx, logger, "nix")
	if err != nil {
		return err
	}

	return s.Shell.Put(ctx, job, os.TempDir())
}
