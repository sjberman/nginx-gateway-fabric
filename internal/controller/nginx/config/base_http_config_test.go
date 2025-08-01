package config

import (
	"sort"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/nginx/nginx-gateway-fabric/v2/internal/controller/state/dataplane"
)

func TestExecuteBaseHttp_HTTP2(t *testing.T) {
	t.Parallel()
	confOn := dataplane.Configuration{
		BaseHTTPConfig: dataplane.BaseHTTPConfig{
			HTTP2: true,
		},
	}

	confOff := dataplane.Configuration{
		BaseHTTPConfig: dataplane.BaseHTTPConfig{
			HTTP2: false,
		},
	}

	expSubStr := "http2 on;"

	tests := []struct {
		name     string
		conf     dataplane.Configuration
		expCount int
	}{
		{
			name:     "http2 on",
			conf:     confOn,
			expCount: 1,
		},
		{
			name:     "http2 off",
			expCount: 0,
			conf:     confOff,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			res := executeBaseHTTPConfig(test.conf)
			g.Expect(res).To(HaveLen(1))
			g.Expect(test.expCount).To(Equal(strings.Count(string(res[0].data), expSubStr)))
			g.Expect(strings.Count(string(res[0].data), "map $http_host $gw_api_compliant_host {")).To(Equal(1))
			g.Expect(strings.Count(string(res[0].data), "map $http_upgrade $connection_upgrade {")).To(Equal(1))
			g.Expect(strings.Count(string(res[0].data), "map $request_uri $request_uri_path {")).To(Equal(1))
		})
	}
}

func TestExecuteBaseHttp_Snippets(t *testing.T) {
	t.Parallel()

	conf := dataplane.Configuration{
		BaseHTTPConfig: dataplane.BaseHTTPConfig{
			Snippets: []dataplane.Snippet{
				{
					Name:     "snippet1",
					Contents: "contents1",
				},
				{
					Name:     "snippet2",
					Contents: "contents2",
				},
			},
		},
	}

	g := NewWithT(t)

	res := executeBaseHTTPConfig(conf)
	g.Expect(res).To(HaveLen(3))

	sort.Slice(
		res, func(i, j int) bool {
			return res[i].dest < res[j].dest
		},
	)

	/*
		Order of files:
		/etc/nginx/conf.d/http.conf
		/etc/nginx/includes/snippet1.conf
		/etc/nginx/includes/snippet2.conf
	*/

	httpRes := string(res[0].data)
	g.Expect(httpRes).To(ContainSubstring("map $http_host $gw_api_compliant_host {"))
	g.Expect(httpRes).To(ContainSubstring("map $http_upgrade $connection_upgrade {"))
	g.Expect(httpRes).To(ContainSubstring("map $request_uri $request_uri_path {"))
	g.Expect(httpRes).To(ContainSubstring("include /etc/nginx/includes/snippet1.conf;"))
	g.Expect(httpRes).To(ContainSubstring("include /etc/nginx/includes/snippet2.conf;"))

	snippet1IncludeRes := string(res[1].data)
	g.Expect(snippet1IncludeRes).To(ContainSubstring("contents1"))

	snippet2IncludeRes := string(res[2].data)
	g.Expect(snippet2IncludeRes).To(ContainSubstring("contents2"))
}

func TestExecuteBaseHttp_NginxReadinessProbePort(t *testing.T) {
	t.Parallel()

	defaultConfig := dataplane.Configuration{
		BaseHTTPConfig: dataplane.BaseHTTPConfig{
			NginxReadinessProbePort: dataplane.DefaultNginxReadinessProbePort,
		},
	}

	customConfig := dataplane.Configuration{
		BaseHTTPConfig: dataplane.BaseHTTPConfig{
			NginxReadinessProbePort: 9090,
		},
	}

	tests := []struct {
		name           string
		expectedPort   string
		expectedListen string
		conf           dataplane.Configuration
	}{
		{
			name:           "default nginx readiness probe port",
			conf:           defaultConfig,
			expectedPort:   "8081",
			expectedListen: "listen 8081;",
		},
		{
			name:           "custom nginx readiness probe 9090",
			conf:           customConfig,
			expectedPort:   "9090",
			expectedListen: "listen 9090;",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			res := executeBaseHTTPConfig(test.conf)
			g.Expect(res).To(HaveLen(1))

			httpConfig := string(res[0].data)

			// check that the listen directive contains the expected port
			g.Expect(httpConfig).To(ContainSubstring(test.expectedListen))

			// check that the health check server block is present
			g.Expect(httpConfig).To(ContainSubstring("server {"))
			g.Expect(httpConfig).To(ContainSubstring("access_log off;"))
			g.Expect(httpConfig).To(ContainSubstring("location = /readyz {"))
			g.Expect(httpConfig).To(ContainSubstring("return 200;"))
		})
	}
}
