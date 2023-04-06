{{- define "postgresDsn" -}}
	{{- $db := .Values.batch.database -}}
	{{- $dsn := printf "host=%s port=%v user=%s password=%s dbname=%s" $db.hostname $db.port $db.user $db.password $db.database -}}
	{{- if $db.tls.enabled -}}
		{{- range $key, $val := $db.tls -}}
			{{- if ne "enabled" $key -}}
				{{- if ne "sslmode" $key -}}
					{{- $val = (include "geoinfo.CertFile" (dict "file" $val)) -}}
				{{- end -}}
				{{- $dsn = printf "%s %s=%s" $dsn $key $val -}}
			{{- end -}}
		{{- end -}}
	{{- else -}}
		{{- $dsn = printf "%s sslmode=disable" $dsn -}}
	{{- end -}}
{{ $dsn }}
{{ end }}

{{- define "postgresDsn2" -}}
    {{- $db := .Values.batch.database -}}
	{{- $pg := .Values.geoinfo.database.credentials -}}
    {{- $dsn := printf "postgres://%s:%s@%s:%v/%s" $db.user $db.password $pg.host $pg.port $pg.dbname -}}
    {{- if $db.tls.enabled -}}
        {{- range $key, $val := $db.tls -}}
            {{- if and (ne "enabled" $key) (ne "sslmode" $key) -}}
                {{- $val = (include "geoinfo.CertFile" (dict "file" $val)) -}}
				{{- $params := "" -}}
				{{- if eq "" $params -}}
					{{- $params = printf "?%s=%s" $key $val -}}
				{{- else -}}
					{{- $params = printf "%s&%s=$s" $params $key $val -}}
				{{- end -}}
                {{- $dsn = printf "%s%s" $dsn $params -}}
            {{- end -}}
        {{- end -}}
        {{- if $db.tls.sslmode -}}
            {{- $dsn = printf "%s&sslmode=%s" $dsn $db.tls.sslmode -}}
        {{- else -}}
            {{- $dsn = printf "%s&sslmode=require" $dsn -}}
        {{- end -}}
    {{- else -}}
        {{- $dsn = printf "%s?sslmode=disable" $dsn -}}
    {{- end -}}
    {{- $dsn -}}
{{- end -}}


{{/*
renderSecret takes the values from .Values.geoinfo.database object to populate
the tom file for database secrets
*/}}
{{- define "renderSecret" -}}
	{{- $secrets := .Values.geoinfo.database.credentials -}}
	{{- $res := "[db]\n" -}}
	{{- range $key, $val := $secrets -}}
		{{- if eq $key "port" -}}
			{{- $res = printf "%s%-15s = %v\n" $res $key $val -}}
		{{- else -}}
			{{- $res = printf "%s%-15s = \"%s\"\n" $res $key $val -}}
		{{- end -}}
	{{- end -}}
	{{- if and .Values.geoinfo.database.tls .Values.geoinfo.database.tls.enable -}}
		{{- range $key, $val := .Values.geoinfo.database.tls -}}
			{{- if contains "ssl" $key -}}
				{{- if ne "sslmode" $key -}}
					{{- $val = (include "geoinfo.CertFile" (dict "file" $val)) -}}
				{{- end -}}
				{{- $res = printf "%s%-15s = \"%s\"\n" $res $key $val -}}
			{{- end -}}
		{{- end -}}
	{{- else -}}
		{{- $res = printf "%s%-15s = \"%s\"\n" $res "sslmode" "disable"}}
	{{- end -}}
	{{- $res -}}
{{ end }}

{{- define "geoinfo.ConfigFile" -}}
{{- $file := required "The name of the config file is required" .file -}}
{{- printf "/etc/geoinfo/%s" $file -}}
{{- end -}}

{{- define "geoinfo.CertFile" -}}
{{- $file := required "The name of the certificate file is required" .file -}}
{{- printf "/etc/geoinfo/creds/%s" $file -}}
{{- end -}}

{{- define "autoGen" -}}
	{{- $ca := genCA "geoinfo" 365 -}}
	{{- $client := genSelfSignedCert "client" (list) (list) 365 -}}
	{{- $cn := .Values.geoinfo.database.credentials.host -}}
	{{- $server := genSelfSignedCert $cn (list "127.0.0.1") (list $cn) 365 -}}
	{{- $certs := dict "ca" $ca "client" $client "server" $server -}}
{{- end }}
