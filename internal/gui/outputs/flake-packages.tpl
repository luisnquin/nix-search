{{ range . }}{{ .package.name }} ({{ .package.version }}){{ if or (gt (len .package.description) 50) (gt (len .package.name) 25)  }}
{{ .package.description }}{{ else }} - {{ .package.description }}{{ end }}

Flake: {{ .flake.name }}
Package: {{ .package.pname }} ({{ .package.set }})
Programs: {{ .package.programs }}
Outputs: {{ .package.outputs }}

{{ if .package.long_description }}{{ .package.long_description }}

{{ end }}{{ if .package.repo_position }}Source: {{ .package.repo_position | transform_source }}
{{ end }}License: {{ .package.license.full_name }}

--
{{ end }}