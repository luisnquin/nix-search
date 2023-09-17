{{ $total_minus_one := sub (len .) 1 }}{{ range $index, $_ := . }}{{ .name }} ({{ .version }}) - {{ .description }}

Package: {{ .pname }} ({{ .set }})
Programs: {{ .programs }}
Outputs: {{ .outputs }}

{{ if .long_description }}{{ .long_description }}

{{ end }}{{ if .repo_position }}Source: {{ .repo_position | transform_source }}
{{ end }}{{ if and (.license) (.license.full_name) }}License: {{ .license.full_name }}{{ end }}

{{ if ne $index $total_minus_one }}--{{ end }}
{{ end }}