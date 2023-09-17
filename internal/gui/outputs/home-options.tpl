{{ range . }}{{ .name }}{{if gt (len .description) 50 }}
{{ .description }}{{ else }} - {{.description}}{{ end }}

Type: {{ .type }}
Example: {{ .example }}
Default: {{ .default }}

{{ if .long_description }}{{ .long_description }}
{{ end }}Source: {{ .source }}

--
{{ end }}