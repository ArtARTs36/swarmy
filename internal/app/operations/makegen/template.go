package makegen

import "github.com/artarts36/swarmy/internal/types/composefile"

const makefileTemplate = `{{ range $stack := .Stacks }}.PHONY: up-{{ $stack.Name }}
up-{{ $stack.Name }}: ## Deploy {{ $stack.Name }} to Swarm
	{{- range $job := $stack.DeployJobs.Before }}
	@echo "> Running job \"{{ $job.Service.Name }}-{{ $job.Name }}\""

	docker run --rm --name {{ $stack.Name }}-{{ $job.Service.Name }}-{{ $job.Name }} --detach=false \
	{{- if $job.PullPolicy }}
	 --pull {{ $job.PullPolicy }} \
	{{- end }}
	{{ if $job.Environment }} 
		{{- range $k, $v := $job.Environment }} --env "{{ $k }}={{ $v }}" {{- end }} \
	{{- end }}
	{{ if $job.Networks }}
		{{- range $network := $job.Networks }} --network {{ $stack.ComposeFile.ResolveNetworkName $network }} {{- end }} \
	{{ else if $job.Service.Networks }}
		{{- range $network := $job.Service.Networks }} --network {{ $stack.ComposeFile.ResolveNetworkName $network }} {{- end }} \
	{{- end }}
	 {{ $job.Image }}
	{{ end }}
	docker stack deploy --with-registry-auth -c {{ $stack.ComposeFilePath }} {{ $stack.Name }} --detach=false

{{ end }}`

type stackSpec struct {
	Name            string
	ComposeFilePath string
	ComposeFile     *composefile.File

	DeployJobs deployJobsSpec
}

type deployJobsSpec struct {
	Before []deployJobSpec
}

type deployJobSpec struct {
	composefile.DeployJob

	Service *composefile.Service
}
