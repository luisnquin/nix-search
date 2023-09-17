{{ range . }}{{ .name }}{{if gt (len .description) 50 }}
{{ .description }}{{ else }} - {{.description}}{{ end }}

Type: {{ .type }}
Example: {{ if .example }}{{ .example }}{{ else }}<nothing>{{ end }}
Default: {{ .default }}

Source: {{if .source }}{{ .source | transform_source }}{{ end }}

--
{{ end }}