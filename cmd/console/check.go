package main

import (
	"context"
	//"errors"
	//"os"

	//"khranity/app/log"
	"khranity/app/lore"
	//"khranity/app/shell"
	"khranity/app/cloud"
)

type CheckCommand struct {
	//Lore  string `short:"l" long:"lore" description:"Lore about the default config"`
	//Name  string `short:"n" long:"name" description:"Name of job to get"`
	//Path  string `short:"p" long:"path" description:"Path to output"`
	//Token string `short:"t" long:"token" description:"Token file to decrypt"`
}

var checkCommand CheckCommand

func (cmd *CheckCommand) Execute(args []string) error {
	// logger 	:= log.GetDefault()
	lore := lore.Load(parser.Groups()[0].FindOptionByLongName("lore").Value().(string))
	// if (lore == nil) {
	// 	return errors.New(" undefined lore file")
	// }

	return cloud.TestClouds(context.TODO(), lore.Clouds)
}

func init() {
	parser.AddCommand("check",
		"Check for access to cloud",
		"The \"check\" command tests access to cloud with <lore> file.",
		&checkCommand)
}