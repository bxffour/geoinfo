{{- define "postgresDsn" -}}
	{{- $db := .Values.batch.database -}}
	{{- $dsn := printf "host=%s port=%v user=%s password=%s dbname=%s" $db.hostname $db.port $db.user $db.password $db.database -}}
	{{- if $db.tlsenabled -}}
		{{- $dsn = printf "%s sslmode=%s sslkey=%s sslcert=%s sslrootcert=%s" $dsn $db.sslmode $db.sslkey $db.sslcert $db.sslrootcert -}}
	{{- else -}}
		{{- $dsn = printf "%s sslmode=disable" $dsn -}}
	{{- end -}}
{{ $dsn }}
{{ end }}

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
		{{- $mountPath := "/etc/geoinfo/certs" -}}
		{{- range $key, $val := .Values.geoinfo.database.tls -}}
			{{- if contains "ssl" $key -}}
				{{- if ne "sslmode" $key -}}
					{{- $val = printf "%s/%s" $mountPath $val -}}
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
{{- $file := default "config.yaml" .Values.geoinfo.configFile -}}
{{- printf "/etc/geoinfo/config/%s" $file -}}
{{- end -}}

{{- define "geoinfo.CredsFile" -}}
{{- $file := default "secret.toml" .Values.geoinfo.credsFile -}}
{{- printf "/etc/geoinfo/config/%s" $file -}}
{{- end -}}

{{- define "autoGen" -}}
	{{- $ca := genCA "geoinfo" 365 -}}
	{{- $client := genSelfSignedCert "client" (list) (list) 365 -}}
	{{- $cn := .Values.geoinfo.database.credentials.host -}}
	{{- $server := genSelfSignedCert $cn (list "127.0.0.1") (list $cn) 365 -}}
	{{- $certs := dict "ca" $ca "client" $client "server" $server -}}
{{- end }}
