package composefile

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type File struct {
	Name string `yaml:"_"`
	Path string `yaml:"-"`

	Services map[string]Service `yaml:"services"`

	Networks map[string]Network `yaml:"networks"`
}

type Network struct {
	Alias string `yaml:"-"`

	Name string `yaml:"name"`
}

type Service struct {
	Name     string   `yaml:"-"`
	Networks []string `yaml:"networks"`

	DeployJobs DeployJobs `yaml:"x-deploy-jobs"`
}

type DeployJobs struct {
	Before []DeployJob `yaml:"before"`
}

type DeployJob struct {
	Name        string            `yaml:"name"`
	PullPolicy  string            `yaml:"pull_policy"` // "always", "missing", "never"
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment"`
	Networks    []string          `yaml:"networks"`
}

func ParseFile(path string) (*File, error) {
	var file File

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	err = yaml.Unmarshal(content, &file)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	file.Name = filepath.Base(path)
	file.Path, err = filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("get abs path: %w", err)
	}

	return &file, nil
}

func (f *File) ResolveNetworkName(network string) string {
	if n, ok := f.Networks[network]; ok {
		if n.Name != "" {
			return n.Name
		}
	}

	return network
}
