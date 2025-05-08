package main

import (
	"github.com/DrSmithFr/go-console/input/option"
	"github.com/artarts36/swarmy/internal/app/operations/makegen"
	"os"

	"github.com/DrSmithFr/go-console"
	"github.com/DrSmithFr/go-console/input/argument"
)

func main() {
	script := go_console.Command{
		Description: "Swarmy",
		Scripts: []*go_console.Script{
			{
				Name:        "make",
				Description: "Generate Makefile of Docker Compose files",
				Arguments: []go_console.Argument{
					{
						Name:        "compose-file",
						Value:       argument.Required | argument.List,
						Description: "paths to docker compose files",
					},
				},
				Options: []go_console.Option{
					{
						Name:        "output",
						Description: "path to output Makefile",
						Value:       option.Optional,
					},
					{
						Name:        "join",
						Description: "path to Makefile, which joined with output",
						Value:       option.Optional,
					},
				},
				Runner: runMakeCommand,
			},
		},
	}

	script.Run()
}

func runMakeCommand(cmd *go_console.Script) go_console.ExitCode {
	output := os.Stdout

	if cmd.Input.Option("output") != "" {
		outputPath := cmd.Input.Option("output")

		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			cmd.Output.Println(err.Error())
			return go_console.ExitError
		}

		output = file
	}

	err := makegen.NewOperation().Run(makegen.Params{
		ComposeFilePaths: cmd.Input.ArgumentList("compose-file"),
		JoinPath:         cmd.Input.Option("join"),
	}, output)
	if err != nil {
		cmd.Output.Println(err.Error())

		return go_console.ExitError
	}
	return go_console.ExitSuccess
}
