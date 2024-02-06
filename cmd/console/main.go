package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	// Example of verbosity
	Verbose bool `short:"v" long:"verbose" description:"Verbose output"`

	Lore string `short:"l" long:"lore" description:"Lore file" default:"lore.yml"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}
