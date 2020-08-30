package main

import (
	"fmt"
	"io"
	"os"
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
	if len(os.Args) < 2 {
		fmt.Printf("usage: tinygo-autocompl --completion-script-bash")
		return nil
	}

	if os.Args[1] == `--completion-script-bash` {
		handleCompletionScriptBash()
	} else {
		fmt.Printf("%s\n", completionBash(os.Args[2:]))
	}

	return nil
}
