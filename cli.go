package main

import (
	"fmt"
	"io"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	appName        = "tinygo-autocmpl"
	appDescription = ""
)

type cli struct {
	outStream io.Writer
	errStream io.Writer
}

var (
	app = kingpin.New(appName, appDescription)
)

// Run ...
func (c *cli) Run(args []string) error {
	app.UsageWriter(c.errStream)

	if VERSION != "" {
		app.Version(fmt.Sprintf("%s version %s build %s", appName, VERSION, BUILDDATE))
	} else {
		app.Version(fmt.Sprintf("%s version - build -", appName))
	}
	app.HelpFlag.Short('h')

	k, err := app.Parse(args[1:])
	if err != nil {
		return err
	}

	switch k {
	default:
	}

	return nil
}
