{{ range . }}{{ .name }}{{ $desc_len := len .description }}{{if gt $desc_len 50 }}
{{ .description }}{{ else }} - {{.description}}{{ end }}

Example: {{ if .example }}{{ .example }}{{ else }}<nothing>{{ end }}
Default: {{ .default }}

Source: {{if .source }}{{ .source | transform_source }}{{ end }}

--
{{ end }}