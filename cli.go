package main

import (
	"bufio"
	"flag"
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
	var (
		completionScriptBash = flag.Bool("completion-script-bash", false, "print completion-script-bash")
		targetsListPath      = flag.String("targets", "", "targets list file")
	)

	flag.Parse()

	if *completionScriptBash {
		handleCompletionScriptBash(*targetsListPath)
		return nil
	}

	if *targetsListPath != "" {
		r, err := os.Open(*targetsListPath)
		if err != nil {
			return err
		}
		defer r.Close()

		scanner := bufio.NewScanner(r)

		targets := []string{}
		for scanner.Scan() {
			targets = append(targets, scanner.Text())
		}
		flagCompleteMap["target"] = targets
	}

	fmt.Printf("%s\n", completionBash(flag.Args()))

	return nil
}
