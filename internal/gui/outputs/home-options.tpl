{{ range . }}{{ .title }} - {{ .description }}
{{if .note }}Note: {{ .note }}
{{ end }}Type: {{ .type }}
Example: {{ .example }}
Default: {{ .default }}
Position: {{ .declared_by }}

{{ end }}