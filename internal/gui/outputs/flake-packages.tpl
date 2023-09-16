{{ range . }}{{ .package.name }} ({{ .package.version }}) - {{ .package.description }}
Flake: {{ .flake.name }}
Programs: {{ .package.programs }}
Outputs: {{ .package.outputs }}

{{ end }}