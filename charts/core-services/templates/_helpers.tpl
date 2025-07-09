{{- define "core-services.name" -}}
{{ .Chart.Name }}
{{- end -}}

{{- define "core-services.fullname" -}}
{{ include "core-services.name" . }}
{{- end -}} 