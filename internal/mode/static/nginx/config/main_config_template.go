package config

const mainConfigTemplateText = `
{{ if .Conf.Telemetry.Endpoint -}}
load_module modules/ngx_otel_module.so;
{{ end -}}

error_log stderr {{ .Conf.Logging.ErrorLevel }};

{{ range $i := .Includes -}}
include {{ $i.Name }};
{{ end -}}
`

const mgmtIncludesTemplateText = `
mgmt {
	{{- if .Endpoint }}
	usage_report endpoint={{ .Endpoint }};
	{{- end }}
	{{- if .Resolver }}
	resolver {{ .Resolver }};
	{{- end }}
	license_token {{ .LicenseTokenFile }};
	{{- if .DeploymentCtxFile }}
	deployment_context {{ .DeploymentCtxFile }};
	{{- end }}
	{{- if .SkipVerify }}
	ssl_verify off;
	{{- end }}
}
`
