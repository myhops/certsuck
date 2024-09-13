{{- range $i, $item := .Verified -}}
{{- println "Verified" $i -}}
	{{- range $j, $val  := $item -}}
		{{- printf "%2d  " $j -}}Subject: {{ .Subject -}}{{- println -}}
		{{- print "    " -}}Issuer:  {{ .Issuer }}{{- println -}}
		{{- print "    " -}}Version:  {{ .Version }}{{- println -}}
	{{- end }}
{{- end -}}

{{- if len .Peer | lt 0 -}}
{{- println "Peer" -}}
	{{- range $j, $val  := .Peer -}}
		{{- printf "%2d  " $j -}}Subject: {{ .Subject -}}{{- println -}}
		{{- print "    " -}}Issuer:  {{ .Issuer }}{{- println -}}
		{{- print "    " -}}Version:  {{ .Version }}{{- println -}}
	{{- end }}
{{- end -}}