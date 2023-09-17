{{ $total_minus_one := sub (len .) 1 }}{{ range $index, $_ := . }}{{ .name }}{{if gt (len .description) 50 }}
{{ .description }}{{ else }} - {{.description}}{{ end }}

Type: {{ .type }}
Example: {{ .example }}
Default: {{ .default }}

{{ if .long_description }}{{ .long_description }}
{{ end }}Source: {{ .source }}

{{ if ne $index $total_minus_one }}--{{ end }}
{{ end }}