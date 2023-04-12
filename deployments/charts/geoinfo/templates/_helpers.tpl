{{- define "postgresDsn" -}}
	{{- $db := .Values.batch.database -}}
	{{- $tls := .Values.geoinfo.database.tls -}}
	{{- $pg := .Values.geoinfo.database.credentials -}}
	{{- $dsn := printf "postgres://%s:%s@%s:%v/%s" $db.user $db.password $pg.host $pg.port $pg.dbname -}}
	{{- if $tls.enable -}}
		{{- $params := ""}}
		{{- range $key, $val := $tls -}}
			{{- if and (contains "ssl" $key) (ne "sslmode" $key) -}}
				{{- if not (and (eq "require" $tls.sslmode) (eq "sslrootcert" $key)) -}}
					{{- $mountPath := "/etc/certs/" -}}
					{{- $val = (include "geoinfo.CertFile" (dict "file" $val "path" $mountPath)) -}}
					{{- if eq "" $params -}}
						{{- $params = printf "?%s=%s" $key $val -}}
					{{- else -}}
						{{- $params = printf "%s&%s=%s" $params $key $val -}}
					{{- end -}}
				{{- end -}}
			{{- end -}}
		{{- end -}}
		{{- $dsn = printf "%s%s" $dsn $params -}}
		{{- if $tls.sslmode -}}
			{{- $dsn = printf "%s&sslmode=%s" $dsn $tls.sslmode -}}
		{{- else -}}
			{{- $dsn = printf "%s&sslmode=prefer" $dsn -}}
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
	{{- $tls := .Values.geoinfo.database.tls -}}
	{{- if and .Values.geoinfo.database.tls .Values.geoinfo.database.tls.enable -}}
		{{- range $key, $val := $tls -}}
			{{- if and (contains "ssl" $key) (ne "sslmode" $key) -}}
				{{- if not (and (eq "require" $tls.sslmode) (eq "sslrootcert" $key)) -}}
					{{- $val = (include "geoinfo.CertFile" (dict "file" $val "path" "/etc/geoinfo/certs/")) -}}
					{{- $res = printf "%s%-15s = \"%s\"\n" $res $key $val -}}
				{{- end -}}
			{{- else if eq "sslmode" $key -}}
				{{- $res = printf "%s%-15s = \"%s\"\n" $res $key $val -}}
			{{- end -}}
		{{- end -}}
	{{- else -}}
		{{- $res = printf "%s%-15s = \"%s\"\n" $res "sslmode" "disable"}}
	{{- end -}}
	{{- $res -}}
{{ end }}

{{- define "geoinfo.ConfigFile" -}}
	{{- required "The name of the config file is required" .file | printf "/etc/geoinfo/%s" -}}
{{- end -}}

{{- define "geoinfo.CertFile" -}}
	{{- $file := required "The name of the certificate file is required" .file  -}}
	{{- $path := required "The mount path for the certificates is required" .path | clean -}}
	{{- printf "%s/%s" $path $file}}
{{- end -}}

{{- define "common.image" -}}
	{{- printf "%s:%s" .name .tag}}
{{- end -}}

{{- define "geoinfo.api.image" -}}
	{{- $image := .Values.geoinfo.image -}}
	{{- include "common.image" (dict "name" $image.name "tag" $image.tag) -}}
{{- end -}}

{{- define "geoinfo.init.image" -}}
	{{- $image := .Values.batch.image -}}
	{{- include "common.image" (dict "name" $image.name "tag" $image.tag) -}}
{{- end -}}

{{- define "geoinfo.createCerts" -}}
{{- $tls := .Values.geoinfo.database.tls -}}
{{- if and $tls.enable $tls.autoGen -}}
	{{- true -}}
{{- end -}}
{{- end -}}