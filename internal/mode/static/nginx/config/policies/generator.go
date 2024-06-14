package policies

import (
	"github.com/nginxinc/nginx-gateway-fabric/internal/mode/static/state/dataplane"
)

type Generator interface {
	GenerateForServer(server dataplane.VirtualServer) GenerateForServerResult
	GenerateForPathRule(rule dataplane.PathRule) GenerateForLocationResult
	GenerateForMatchRule(rule dataplane.MatchRule) GenerateForLocationResult
}

type GenerateForServerResult []ServerFile

type ServerFile struct {
	Name    string
	Content []byte
}

type GenerateForLocationResult []LocationFile

type LocationFile struct {
	Name    string
	Type    LocationType
	Content []byte
}

type LocationType int

const (
	Internal LocationType = iota
	External
	ExternalRedirect
)

func (locFile GenerateForLocationResult) ForInternalLocation() []LocationFile {
	files := make([]LocationFile, 0, len(locFile))
	for _, f := range locFile {
		if f.Type == Internal {
			files = append(files, f)
		}
	}

	return files
}

func (locFile GenerateForLocationResult) ForExternalLocation() []LocationFile {
	files := make([]LocationFile, 0, len(locFile))
	for _, f := range locFile {
		if f.Type == External {
			files = append(files, f)
		}
	}

	return files
}

func (locFile GenerateForLocationResult) ForExternalRedirectLocation() []LocationFile {
	files := make([]LocationFile, 0, len(locFile))
	for _, f := range locFile {
		if f.Type == ExternalRedirect {
			files = append(files, f)
		}
	}

	return files
}

type CompositeGenerator struct {
	generators []Generator
}

func NewCompositeGenerator(generators ...Generator) *CompositeGenerator {
	return &CompositeGenerator{generators: generators}
}

func (g *CompositeGenerator) GenerateForServer(server dataplane.VirtualServer) GenerateForServerResult {
	var results GenerateForServerResult

	for _, generator := range g.generators {
		results = append(results, generator.GenerateForServer(server)...)
	}

	return results
}

func (g *CompositeGenerator) GenerateForPathRule(rule dataplane.PathRule) GenerateForLocationResult {
	var results GenerateForLocationResult

	for _, generator := range g.generators {
		results = append(results, generator.GenerateForPathRule(rule)...)
	}

	return results
}

func (g *CompositeGenerator) GenerateForMatchRule(rule dataplane.MatchRule) GenerateForLocationResult {
	var results GenerateForLocationResult

	for _, generator := range g.generators {
		results = append(results, generator.GenerateForMatchRule(rule)...)
	}

	return results
}
