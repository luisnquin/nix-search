{{ $total_minus_one := sub (len .) 1 }}{{ range $index, $_ := . }}{{ .option.name }}{{ if gt (len .option.description) 50 }}
{{ .option.description }}{{ else }} - {{ .option.description }}{{ end }}

Flake: {{ .flake.name }}
Type: {{ .option.type }}
Example: {{ .option.example }}
Default: {{ .option.default }}
{{ if .option.source }}
Source: {{ .option.source | transform_source }}{{ end }}

{{ if ne $index $total_minus_one }}--{{ end }}
{{ end }}