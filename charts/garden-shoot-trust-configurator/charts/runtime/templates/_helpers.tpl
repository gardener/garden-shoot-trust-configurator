{{- define "garden-shoot-trust-configurator.name" -}}
garden-shoot-trust-configurator
{{- end -}}

{{- define "image" -}}
    {{- if hasPrefix "sha256:" (required "$.tag is required" $.tag) -}}
        {{ required "$.repository is required" $.repository }}@{{ required "$.tag is required" $.tag }}
    {{- else -}}
        {{ required "$.repository is required" $.repository }}:{{ required "$.tag is required" $.tag }}
    {{- end -}}
{{- end -}}

{{- define "garden-shoot-trust-configurator.config.data" -}}
config.yaml: |
{{ include "garden-shoot-trust-configurator.config" . | indent 2 }}
{{- end -}}

{{- define "garden-shoot-trust-configurator.config.name" -}}
garden-shoot-trust-configurator-configmap
{{- end -}}


{{- define "garden-shoot-trust-configurator.config" -}}
apiVersion: config.trust-configurator.gardener.cloud/v1alpha1
kind: GardenShootTrustConfiguratorConfiguration
logLevel: {{ .Values.config.logLevel }}
logFormat: {{ .Values.config.logFormat }}
controllers:
  shoot:
    syncPeriod: {{ .Values.config.controllers.shoot.syncPeriod }}
    oidcConfig:
      maxTokenExpiration: {{ .Values.config.controllers.shoot.oidcConfig.maxTokenExpiration }}
      audiences:
{{ toYaml .Values.config.controllers.shoot.oidcConfig.audiences | indent 6 }}
  garbageCollector:
    syncPeriod: {{  .Values.config.controllers.garbageCollector.syncPeriod }}
    minimumObjectLifetime: {{  .Values.config.controllers.garbageCollector.minimumObjectLifetime }}
server:
  healthPort: {{ .Values.config.server.healthPort }}
  port: {{ .Values.config.server.port }}
leaderElection:
  {{- if .Values.config.leaderElection.leaderElect }}
  leaderElect: {{ .Values.config.leaderElection.leaderElect }}
  {{- end }}
  {{- if .Values.config.leaderElection.leaseDuration }}
  leaseDuration: {{ .Values.config.leaderElection.leaseDuration }}
  {{- end }}
  {{- if .Values.config.leaderElection.renewDeadline }}
  renewDeadline: {{ .Values.config.leaderElection.renewDeadline }}
  {{- end }}
  {{- if .Values.config.leaderElection.retryPeriod }}
  retryPeriod: {{ .Values.config.leaderElection.retryPeriod }}
  {{- end }}
  {{- if .Values.config.leaderElection.resourceLock }}
  resourceLock: {{ .Values.config.leaderElection.resourceLock }}
  {{- end }}
  {{- if .Values.config.leaderElection.resourceName }}
  resourceName: {{ .Values.config.leaderElection.resourceName }}
  {{- end }}
  {{- if .Values.config.leaderElection.resourceNamespace }}
  resourceNamespace: {{ .Values.config.leaderElection.resourceNamespace }}
  {{- end }}
{{- end -}}
