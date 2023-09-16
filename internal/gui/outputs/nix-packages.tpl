{{range .}}{{.name}} ({{.version}}) - {{.description}}
{{ if .long_description }}Note: {{.long_description}}\n{{end}}Programs: {{.programs}}
Outputs: {{.outputs}}
Source: {{.source}}

{{end}}