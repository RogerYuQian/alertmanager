package cli

import (
	"fmt"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/template"
	"gopkg.in/alecthomas/kingpin.v2"
)

// TODO: This can just be a type that is []string, doesn't have to be a struct
type checkConfigCmd struct {
	files []string
}

const checkConfigHelp = `Validate alertmanager config files

Will validate the syntax and schema for alertmanager config file
and associated templates. Non existing templates will not trigger
errors.
`

func configureCheckConfigCmd(app *kingpin.Application) {
	var (
		c        = &checkConfigCmd{}
		checkCmd = app.Command("check-config", checkConfigHelp)
	)
	checkCmd.Arg("check-files", "Files to be validated").ExistingFilesVar(&c.files)
	checkCmd.Action(c.checkConfig)
}

func (c *checkConfigCmd) checkConfig(ctx *kingpin.ParseContext) error {
	return CheckConfig(c.files)
}

func CheckConfig(args []string) error {
	failed := 0

	for _, arg := range args {
		fmt.Printf("Checking '%s'", arg)
		config, _, err := config.LoadFile(arg)
		if err != nil {
			fmt.Printf("  FAILED: %s\n", err)
			failed++
		} else {
			fmt.Printf("  SUCCESS\n")
		}

		if config != nil {
			fmt.Printf("Found %d templates: ", len(config.Templates))
			if len(config.Templates) > 0 {
				_, err = template.FromGlobs(config.Templates...)
				if err != nil {
					fmt.Printf("  FAILED: %s\n", err)
					failed++
				} else {
					fmt.Printf("  SUCCESS\n")
				}
			}
		}
		fmt.Printf("\n")
	}
	if failed > 0 {
		return fmt.Errorf("failed to validate %d file(s)", failed)
	}
	return nil
}
