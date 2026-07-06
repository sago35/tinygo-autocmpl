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
		completionScriptBash  = flag.Bool("completion-script-bash", false, "print completion-script-bash")
		completionScriptZsh   = flag.Bool("completion-script-zsh", false, "print completion-script-zsh")
		completionScriptClink = flag.Bool("completion-script-clink", false, "print completion-script-clink")
		completionScriptFish  = flag.Bool("completion-script-fish", false, "print completion-script-fish")
		targetsListPath       = flag.String("targets", "", "targets list file")
		showVersion           = flag.Bool("version", false, "print version information")
	)

	flag.Parse()

	if *showVersion {
		version := VERSION
		if version == "" {
			version = "dev"
		}
		if BUILDDATE != "" {
			fmt.Fprintf(c.outStream, "%s version %s (%s)\n", appName, version, BUILDDATE)
		} else {
			fmt.Fprintf(c.outStream, "%s version %s\n", appName, version)
		}
		return nil
	}

	if *completionScriptBash {
		handleCompletionScriptBash(*targetsListPath)
		return nil
	}

	if *completionScriptZsh {
		handleCompletionScriptZsh(*targetsListPath)
		return nil
	}

	if *completionScriptClink {
		handleCompletionScriptClink(*targetsListPath)
		return nil
	}

	if *completionScriptFish {
		handleCompletionScriptFish(*targetsListPath)
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

	fmt.Printf("%s\n", completeArgs(flag.Args()))

	return nil
}
