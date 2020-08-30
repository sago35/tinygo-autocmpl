package main

import (
	"flag"
	"fmt"
	"io"
)

const (
	appName        = "tinygo-autocmpl"
	appDescription = ""
)

type cli struct {
	outStream io.Writer
	errStream io.Writer
}

// Run ...
func (c *cli) Run(args []string) error {
	var (
		completionScriptBash = flag.Bool("completion-script-bash", false, "print completion-script-bash")
	)

	flag.Parse()

	if *completionScriptBash {
		handleCompletionScriptBash()
		return nil
	}

	fmt.Printf("%s\n", completionBash(flag.Args()))

	return nil
}
