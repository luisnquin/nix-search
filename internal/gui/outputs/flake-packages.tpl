{{ range . }}{{ .package.name }} ({{ .package.version }}){{ $name_len := len .package.name }}{{ $desc_len := len .package.description }}{{ if  or (gt $desc_len 50) (gt $name_len 25)  }}
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