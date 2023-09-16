{{ range . }}{{ .title }} - {{ .description }}

{{if .note }}Note: {{ .note }}

{{ end }}Type: {{ .type }}
Example: {{ .example }}
Default: {{ .default }}

Source: {{ .declared_by }}

--
{{ end }}