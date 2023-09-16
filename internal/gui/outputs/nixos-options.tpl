{{ range . }}{{ .name }} - {{ .description}}
Example: {{ if .example }}{{ .example }}{{ else }}<nothing>{{ end }}
Default: {{ .default }}
Source: {{if .source }}{{ .source | transform_source }}{{ end }}

{{ end }}