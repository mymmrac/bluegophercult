{{- $resp := getPage "http://example.com" -}}
{{- $body := readAll $resp -}}
{{- $body -}}
