{{ $total_minus_one := sub (len .) 1 }}{{ range $index, $_ := . }}{{ .name }}{{if gt (len .description) 50 }}
{{ .description }}{{ else }} - {{.description}}{{ end }}

Type: {{ .type }}
Example: {{ if .example }}{{ .example }}{{ else }}<nothing>{{ end }}
Default: {{ .default }}

Source: {{if .source }}{{ .source | transform_source }}{{ end }}

{{ if ne $index $total_minus_one }}--{{ end }}
{{ end }}