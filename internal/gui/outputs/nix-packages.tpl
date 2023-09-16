{{ range . }}{{ .name }} ({{ .version }}) - {{ .description }}
{{ if .long_description }}Note: {{ .long_description }}
{{ end }}Programs: {{ .programs }}
Outputs: {{ .outputs }}
{{ if .repo_position }}Source: {{ .repo_position | transform_source }}{{ end }}

{{ end }}