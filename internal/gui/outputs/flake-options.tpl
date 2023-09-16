{{ range . }}{{ .option.name }} - {{ .option.description }}
Flake: {{ .flake.name }}
Example: {{ .option.example }}
Default: {{ .option.default }}

{{ end }}