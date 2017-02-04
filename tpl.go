package main

import (
	"bytes"
	"fmt"
	"io"
	log "github.com/Sirupsen/logrus"
	"os"
	"text/template"
	"github.com/urfave/cli"
	"strings"
)

func handleDefault(defaultValue string, providedValue string) string {
	return providedValue
}

func main() {
	log.SetLevel(log.WarnLevel)
	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringSliceFlag {
			Name: "set",
			Usage: "Set a template value (--set name=value)",
		},
		cli.BoolFlag {
			Name: "debug",
			Usage: "Debug output",
		},
	}

	app.Action = func(c *cli.Context) error {
		filePath := ""
		debug := c.Bool("debug")
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		if c.NArg() > 0 {
			filePath = c.Args().Get(0)
		} else {
			log.Fatal("Missing arg: filePath")
		}

		funcs := template.FuncMap {
			"default": handleDefault,
		}
		templ, err := template.New("cs.yaml").Funcs(funcs).ParseFiles(filePath)

		if err != nil {
			return fmt.Errorf("Error parsing template(s): %v", err)
		}
		valuesMap := map[string]string {}

		setVariables := c.StringSlice("set")
		log.Debug("Got params: ", len(setVariables))

		for _, setVariable := range setVariables {
			nameVal := strings.Split(setVariable, "=")
			if (len(nameVal) != 2) {
				log.Fatal("Invalid --set flag: ", nameVal)
			}
			valuesMap[nameVal[0]] = nameVal[1]
			log.Debug("Substituting template: ", nameVal[0], " with value: ", nameVal[1])
		}

		finalValues := map[string]interface{} {
			"Values": valuesMap,
		}
		if true {
			var targetWriter io.Writer = os.Stdout
			if debug {
				targetWriter = bytes.NewBufferString("")
			}

			err = templ.Execute(targetWriter, finalValues)
			if err != nil {
				log.Error("Failed to execute template: ", err)
			} else {
				log.Debug("Executed correctly.")
			}
		}

		return nil
	}

	app.Run(os.Args)
}
