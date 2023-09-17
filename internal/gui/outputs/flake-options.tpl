{{ range . }}{{ .option.name }}{{ if gt (len .option.description) 50 }}
{{ .option.description }}{{ else }} - {{ .option.description }}{{ end }}

Flake: {{ .flake.name }}
Type: {{ .option.type }}
Example: {{ .option.example }}
Default: {{ .option.default }}
{{ if .option.source }}
Source: {{ .option.source | transform_source }}{{ end }}
--
{{ end }}