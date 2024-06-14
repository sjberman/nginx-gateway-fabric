package clientsettings

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	ngfAPI "github.com/nginxinc/nginx-gateway-fabric/apis/v1alpha1"
	"github.com/nginxinc/nginx-gateway-fabric/internal/framework/helpers"
	"github.com/nginxinc/nginx-gateway-fabric/internal/mode/static/nginx/config/policies"
	"github.com/nginxinc/nginx-gateway-fabric/internal/mode/static/state/dataplane"
)

var (
	tmpl                 = template.Must(template.New("client settings policy").Parse(clientSettingsTemplate))
	tmplExternalRedirect = template.Must(
		template.New("client settings policy ext redirect").Parse(externalRedirectTemplate),
	)
)

const clientSettingsTemplate = `
{{- if .Body }}
	{{- if .Body.MaxSize }}
client_max_body_size {{ .Body.MaxSize }};
	{{- end }}
	{{- if .Body.Timeout }}
client_body_timeout {{ .Body.Timeout }};
	{{- end }}
{{- end }}
{{- if .KeepAlive }}
	{{- if .KeepAlive.Requests }}
keepalive_requests {{ .KeepAlive.Requests }};
	{{- end }}
	{{- if .KeepAlive.Time }}
keepalive_time {{ .KeepAlive.Time }};
	{{- end }}
    {{- if .KeepAlive.Timeout }}
        {{- if and .KeepAlive.Timeout.Server .KeepAlive.Timeout.Header }}
keepalive_timeout {{ .KeepAlive.Timeout.Server }} {{ .KeepAlive.Timeout.Header }};
        {{- else if .KeepAlive.Timeout.Server }}
keepalive_timeout {{ .KeepAlive.Timeout.Server }};
        {{- end }}
    {{- end }}
{{- end }}
`

const externalRedirectTemplate = `
client_max_body_size {{ . }};
`

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g Generator) GenerateForServer(server dataplane.VirtualServer) policies.GenerateForServerResult {
	results := make(policies.GenerateForServerResult, 0, len(server.Policies))

	for _, pol := range server.Policies {
		switch csp := pol.(type) {
		case *ngfAPI.ClientSettingsPolicy:
			results = append(results, policies.ServerFile{
				Name:    fmt.Sprintf("%s_%s_%s.conf", csp.Kind, csp.Namespace, csp.Name),
				Content: helpers.MustExecuteTemplate(tmpl, csp.Spec),
			})
		}
	}

	return results
}

func (g Generator) GenerateForPathRule(rule dataplane.PathRule) policies.GenerateForLocationResult {
	var maxBodySize ngfAPI.Size

	for _, match := range rule.MatchRules {
		for _, pol := range match.Policies {
			csp, ok := pol.(*ngfAPI.ClientSettingsPolicy)
			if !ok {
				continue
			}

			if csp.Spec.Body != nil {
				maxBodySize = getMaxSize(maxBodySize, csp.Spec.Body.MaxSize)
			}
		}
	}

	if maxBodySize == "" {
		return nil
	}

	return policies.GenerateForLocationResult{
		{
			Name: fmt.Sprintf(
				"ClientSettingsPolicy_%s_%s_redirect.conf",
				rule.PathType,
				removeSlashFromPath(rule.Path),
			),
			Type:    policies.ExternalRedirect,
			Content: helpers.MustExecuteTemplate(tmplExternalRedirect, maxBodySize),
		},
	}
}

func removeSlashFromPath(path string) string {
	return strings.ReplaceAll(path, "/", "")
}

func (g Generator) GenerateForMatchRule(rule dataplane.MatchRule) policies.GenerateForLocationResult {
	results := make(policies.GenerateForLocationResult, 0, len(rule.Policies)*2)

	for _, pol := range rule.Policies {
		switch csp := pol.(type) {
		case *ngfAPI.ClientSettingsPolicy:
			fileContents := helpers.MustExecuteTemplate(tmpl, csp.Spec)

			results = append(
				results,
				policies.LocationFile{
					Name:    fmt.Sprintf("%s_%s_%s_int.conf", csp.Kind, csp.Namespace, csp.Name),
					Content: fileContents,
					Type:    policies.Internal,
				},
				policies.LocationFile{
					Name:    fmt.Sprintf("%s_%s_%s_ext.conf", csp.Kind, csp.Namespace, csp.Name),
					Content: fileContents,
					Type:    policies.External,
				},
			)
		}
	}

	return results
}

func isDigit(char string) bool {
	_, err := strconv.Atoi(char)
	return err == nil
}

// ^\d{1,4}(k|m|g)?$`
func getMaxSize(s1 ngfAPI.Size, s2 *ngfAPI.Size) ngfAPI.Size {
	if s2 == nil {
		return s1
	}

	if s1 == "" {
		return *s2
	}

	s1Str := string(s1)
	s2Str := string(*s2)

	s1Unit := s1Str[len(s1Str)-1:]
	s2Unit := s2Str[len(s2Str)-1:]

	// if unit is missing then it's bytes
	// bytes will always be smaller than other units
	if isDigit(s1Unit) && !isDigit(s2Unit) {
		return *s2
	}

	if !isDigit(s1Unit) && isDigit(s2Unit) {
		return s1
	}

	if s1Unit == s2Unit {
		s1Int, err := strconv.Atoi(s1Str[:len(s1Str)-1])
		if err != nil {
			panic(err)
		}

		s2Int, err := strconv.Atoi(s2Str[:len(s2Str)-1])
		if err != nil {
			panic(err)
		}

		if s1Int > s2Int {
			return s1
		}

		return *s2
	}

	switch s1Unit {
	case "k":
		return *s2
	case "m":
		if s2Unit == "g" {
			return *s2
		}

		return s1
	case "g":
		return s1
	}

	return s1
}
