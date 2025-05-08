package makegen

import (
	"fmt"
	"html/template"
	"io"

	"github.com/artarts36/swarmy/internal/shared/fpath"
	"github.com/artarts36/swarmy/internal/types/composefile"
)

type Operation struct {
}

type Params struct {
	ComposeFilePaths []string
}

func NewOperation() *Operation {
	return &Operation{}
}

func (op *Operation) Run(params Params, result io.Writer) error {
	stacks := []*stackSpec{}

	for _, path := range params.ComposeFilePaths {
		stack, err := op.parseStack(path)
		if err != nil {
			return fmt.Errorf("parse stack: %w", err)
		}

		stacks = append(stacks, stack)
	}

	tmpl, err := template.New("makefile").Parse(makefileTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	err = tmpl.Execute(result, map[string]interface{}{
		"Stacks": stacks,
	})
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	return nil
}

func (op *Operation) parseStack(path string) (*stackSpec, error) {
	file, err := composefile.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("parse compose file: %w", err)
	}

	stack := &stackSpec{
		Name:            fpath.OmitExt(file.Name),
		ComposeFilePath: file.Path,
		ComposeFile:     file,
	}

	for _, service := range file.Services {
		for _, job := range service.DeployJobs.Before {
			stack.DeployJobs.Before = append(stack.DeployJobs.Before, deployJobSpec{
				DeployJob: job,
				Service:   service,
			})
		}
	}

	return stack, nil
}
