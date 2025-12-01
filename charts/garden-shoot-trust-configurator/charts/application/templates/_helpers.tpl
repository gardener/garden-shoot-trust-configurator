{{- define "garden-shoot-trust-configurator.name" -}}
garden-shoot-trust-configurator
{{- end -}}

# Warning: The following helper is duplicated in charts/runtime/templates/_helpers.tpl. Keep them in sync.
{{- define "leaderelectionid" -}}
garden-shoot-trust-configurator-leader-election
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
