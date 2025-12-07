{{- define "garden-shoot-trust-configurator.name" -}}
garden-shoot-trust-configurator
{{- end -}}

{{- define "leaderelectionid" -}}
{{ .Values.global.leaderElection.id }}
{{- end -}}

{{- define "labels.app.key" -}}
app.kubernetes.io/name
{{- end -}}
{{- define "labels.app.value" -}}
{{ include "garden-shoot-trust-configurator.name" . }}
{{- end -}}

{{- define "labels" -}}
{{ include "labels.app.key" . }}: {{ include "labels.app.value" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
