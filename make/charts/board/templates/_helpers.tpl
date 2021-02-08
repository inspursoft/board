{{- define "board.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "board.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "board.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* board apiserver name, fullname, repository, tag */}}
{{/* we use 60-20-43 lenght limit because we will use the name append 'env','kubeconfig' etc. */}}
{{- define "board.apiserver.name" -}}
{{- $apiserver := default "apiserver" .Values.apiserver.name -}}
{{- $lenght := sub 43 (len $apiserver) -}}
{{- $name := (include "board.name" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $apiserver -}}
{{- end -}}

{{- define "board.apiserver.fullname" -}}
{{- $apiserver := default "apiserver" .Values.apiserver.name -}}
{{- $lenght := sub 43 (len $apiserver) -}}
{{- $name := (include "board.fullname" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $apiserver -}}
{{- end -}}

{{- define "board.apiserver.image.repository" -}}
{{- if .Values.image.registry -}}
{{- printf "%s/%s" .Values.image.registry .Values.apiserver.image.repository -}}
{{- else -}}
{{- .Values.apiserver.image.repository -}}
{{- end -}}
{{- end -}}

{{- define "board.apiserver.image.tag" -}}
{{- default .Values.apiserver.image.tag .Values.image.tag -}}
{{- end -}}

{{- define "board.apiserver.image.image" -}}
{{- printf "\"%s:%s\"" (include "board.apiserver.image.repository" .) (include "board.apiserver.image.tag" .) -}}
{{- end -}}

{{/* board db name, fullname, repository, tag */}}
{{/* we use 60-20-43 lenght limit because we will use the name append 'env','kubeconfig' etc. */}}
{{- define "board.db.name" -}}
{{- $db := default "db" .Values.db.name -}}
{{- $lenght := sub 43 (len $db) -}}
{{- $name := (include "board.name" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $db -}}
{{- end -}}

{{- define "board.db.fullname" -}}
{{- $db := default "db" .Values.db.name -}}
{{- $lenght := sub 43 (len $db) -}}
{{- $name := (include "board.fullname" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $db -}}
{{- end -}}

{{- define "board.db.image.repository" -}}
{{- if .Values.image.registry -}}
{{- printf "%s/%s" .Values.image.registry .Values.db.image.repository -}}
{{- else -}}
{{- .Values.db.image.repository -}}
{{- end -}}
{{- end -}}

{{- define "board.db.image.tag" -}}
{{- default .Values.db.image.tag .Values.image.tag -}}
{{- end -}}

{{- define "board.db.image.image" -}}
{{- printf "\"%s:%s\"" (include "board.db.image.repository" .) (include "board.db.image.tag" .) -}}
{{- end -}}


{{/* board tokenserver name, fullname, repository, tag */}}
{{/* we use 60-20-43 lenght limit because we will use the name append 'env','kubeconfig' etc. */}}
{{- define "board.tokenserver.name" -}}
{{- $tokenserver := default "tokenserver" .Values.tokenserver.name -}}
{{- $lenght := sub 43 (len $tokenserver) -}}
{{- $name := (include "board.name" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $tokenserver -}}
{{- end -}}

{{- define "board.tokenserver.fullname" -}}
{{- $tokenserver := default "tokenserver" .Values.tokenserver.name -}}
{{- $lenght := sub 43 (len $tokenserver) -}}
{{- $name := (include "board.fullname" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $tokenserver -}}
{{- end -}}

{{- define "board.tokenserver.image.repository" -}}
{{- if .Values.image.registry -}}
{{- printf "%s/%s" .Values.image.registry .Values.tokenserver.image.repository -}}
{{- else -}}
{{- .Values.tokenserver.image.repository -}}
{{- end -}}
{{- end -}}

{{- define "board.tokenserver.image.tag" -}}
{{- default .Values.tokenserver.image.tag .Values.image.tag -}}
{{- end -}}

{{- define "board.tokenserver.image.image" -}}
{{- printf "\"%s:%s\"" (include "board.tokenserver.image.repository" .) (include "board.tokenserver.image.tag" .) -}}
{{- end -}}


{{/* board proxy name, fullname, repository, tag */}}
{{/* we use 60-20-43 lenght limit because we will use the name append 'env','kubeconfig' etc. */}}
{{- define "board.proxy.name" -}}
{{- $proxy := default "proxy" .Values.proxy.name -}}
{{- $lenght := sub 43 (len $proxy) -}}
{{- $name := (include "board.name" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $proxy -}}
{{- end -}}

{{- define "board.proxy.fullname" -}}
{{- $proxy := default "proxy" .Values.proxy.name -}}
{{- $lenght := sub 43 (len $proxy) -}}
{{- $name := (include "board.fullname" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $proxy -}}
{{- end -}}

{{- define "board.proxy.image.repository" -}}
{{- if .Values.image.registry -}}
{{- printf "%s/%s" .Values.image.registry .Values.proxy.image.repository -}}
{{- else -}}
{{- .Values.proxy.image.repository -}}
{{- end -}}
{{- end -}}

{{- define "board.proxy.image.tag" -}}
{{- default .Values.proxy.image.tag .Values.image.tag -}}
{{- end -}}

{{- define "board.proxy.image.image" -}}
{{- printf "\"%s:%s\"" (include "board.proxy.image.repository" .) (include "board.proxy.image.tag" .) -}}
{{- end -}}


{{/* board chartmuseum name, fullname, repository, tag */}}
{{/* we use 60-20-43 lenght limit because we will use the name append 'env','kubeconfig' etc. */}}
{{- define "board.chartmuseum.name" -}}
{{- $chartmuseum := default "chartmuseum" .Values.chartmuseum.name -}}
{{- $lenght := sub 43 (len $chartmuseum) -}}
{{- $name := (include "board.name" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $chartmuseum -}}
{{- end -}}

{{- define "board.chartmuseum.fullname" -}}
{{- $chartmuseum := default "chartmuseum" .Values.chartmuseum.name -}}
{{- $lenght := sub 43 (len $chartmuseum) -}}
{{- $name := (include "board.fullname" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $chartmuseum -}}
{{- end -}}

{{- define "board.chartmuseum.image.repository" -}}
{{- if .Values.image.registry -}}
{{- printf "%s/%s" .Values.image.registry .Values.chartmuseum.image.repository -}}
{{- else -}}
{{- .Values.chartmuseum.image.repository -}}
{{- end -}}
{{- end -}}

{{- define "board.chartmuseum.image.tag" -}}
{{- default .Values.chartmuseum.image.tag .Values.image.tag -}}
{{- end -}}

{{- define "board.chartmuseum.image.image" -}}
{{- printf "\"%s:%s\"" (include "board.chartmuseum.image.repository" .) (include "board.chartmuseum.image.tag" .) -}}
{{- end -}}


{{/* board prometheus name, fullname, repository, tag */}}
{{/* we use 60-20-43 lenght limit because we will use the name append 'env','kubeconfig' etc. */}}
{{- define "board.prometheus.name" -}}
{{- $prometheus := default "prometheus" .Values.prometheus.name -}}
{{- $lenght := sub 43 (len $prometheus) -}}
{{- $name := (include "board.name" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $prometheus -}}
{{- end -}}

{{- define "board.prometheus.fullname" -}}
{{- $prometheus := default "prometheus" .Values.prometheus.name -}}
{{- $lenght := sub 43 (len $prometheus) -}}
{{- $name := (include "board.fullname" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $prometheus -}}
{{- end -}}

{{- define "board.prometheus.image.repository" -}}
{{- if .Values.image.registry -}}
{{- printf "%s/%s" .Values.image.registry .Values.prometheus.image.repository -}}
{{- else -}}
{{- .Values.prometheus.image.repository -}}
{{- end -}}
{{- end -}}

{{- define "board.prometheus.image.tag" -}}
{{- default .Values.prometheus.image.tag .Values.image.tag -}}
{{- end -}}

{{- define "board.prometheus.image.image" -}}
{{- printf "\"%s:%s\"" (include "board.prometheus.image.repository" .) (include "board.prometheus.image.tag" .) -}}
{{- end -}}



{{/* board grafana name, fullname, repository, tag */}}
{{/* we use 60-20-43 lenght limit because we will use the name append 'env','kubeconfig' etc. */}}
{{- define "board.grafana.name" -}}
{{- $grafana := default "grafana" .Values.grafana.name -}}
{{- $lenght := sub 43 (len $grafana) -}}
{{- $name := (include "board.name" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $grafana -}}
{{- end -}}

{{- define "board.grafana.fullname" -}}
{{- $grafana := default "grafana" .Values.grafana.name -}}
{{- $lenght := sub 43 (len $grafana) -}}
{{- $name := (include "board.fullname" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $grafana -}}
{{- end -}}

{{- define "board.grafana.image.repository" -}}
{{- if .Values.image.registry -}}
{{- printf "%s/%s" .Values.image.registry .Values.grafana.image.repository -}}
{{- else -}}
{{- .Values.grafana.image.repository -}}
{{- end -}}
{{- end -}}

{{- define "board.grafana.image.tag" -}}
{{- default .Values.grafana.image.tag .Values.image.tag -}}
{{- end -}}

{{- define "board.grafana.image.image" -}}
{{- printf "\"%s:%s\"" (include "board.grafana.image.repository" .) (include "board.grafana.image.tag" .) -}}
{{- end -}}


{{/* board elasticsearch name, fullname, repository, tag */}}
{{/* we use 60-20-43 lenght limit because we will use the name append 'env','kubeconfig' etc. */}}
{{- define "board.elasticsearch.name" -}}
{{- $elasticsearch := default "elasticsearch" .Values.elasticsearch.name -}}
{{- $lenght := sub 43 (len $elasticsearch) -}}
{{- $name := (include "board.name" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $elasticsearch -}}
{{- end -}}

{{- define "board.elasticsearch.fullname" -}}
{{- $elasticsearch := default "elasticsearch" .Values.elasticsearch.name -}}
{{- $lenght := sub 43 (len $elasticsearch) -}}
{{- $name := (include "board.fullname" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $elasticsearch -}}
{{- end -}}

{{- define "board.elasticsearch.image.repository" -}}
{{- if .Values.image.registry -}}
{{- printf "%s/%s" .Values.image.registry .Values.elasticsearch.image.repository -}}
{{- else -}}
{{- .Values.elasticsearch.image.repository -}}
{{- end -}}
{{- end -}}

{{- define "board.elasticsearch.image.tag" -}}
{{- default .Values.elasticsearch.image.tag .Values.image.tag -}}
{{- end -}}

{{- define "board.elasticsearch.image.image" -}}
{{- printf "\"%s:%s\"" (include "board.elasticsearch.image.repository" .) (include "board.elasticsearch.image.tag" .) -}}
{{- end -}}


{{/* board kibana name, fullname, repository, tag */}}
{{/* we use 60-20-43 lenght limit because we will use the name append 'env','kubeconfig' etc. */}}
{{- define "board.kibana.name" -}}
{{- $kibana := default "kibana" .Values.kibana.name -}}
{{- $lenght := sub 43 (len $kibana) -}}
{{- $name := (include "board.name" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $kibana -}}
{{- end -}}

{{- define "board.kibana.fullname" -}}
{{- $kibana := default "kibana" .Values.kibana.name -}}
{{- $lenght := sub 43 (len $kibana) -}}
{{- $name := (include "board.fullname" . | trunc (int $lenght)) -}}
{{- printf "%s-%s" $name $kibana -}}
{{- end -}}

{{- define "board.kibana.image.repository" -}}
{{- if .Values.image.registry -}}
{{- printf "%s/%s" .Values.image.registry .Values.kibana.image.repository -}}
{{- else -}}
{{- .Values.kibana.image.repository -}}
{{- end -}}
{{- end -}}

{{- define "board.kibana.image.tag" -}}
{{- default .Values.kibana.image.tag .Values.image.tag -}}
{{- end -}}

{{- define "board.kibana.image.image" -}}
{{- printf "\"%s:%s\"" (include "board.kibana.image.repository" .) (include "board.kibana.image.tag" .) -}}
{{- end -}}