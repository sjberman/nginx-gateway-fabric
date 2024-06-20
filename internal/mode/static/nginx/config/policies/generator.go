package policies

import (
	"github.com/nginxinc/nginx-gateway-fabric/internal/mode/static/nginx/config/http"
)

type Generator interface {
	GenerateForLocation(policies []Policy, location http.Location) GenerateResultFiles
	// TODO: do we need this as a separate method?
	GenerateForInternalLocation(policies []Policy) GenerateResultFiles
	GenerateForServer(policies []Policy, server http.Server) GenerateResultFiles
}

type GenerateResultFiles []File

type File struct {
	Name    string
	Content []byte
}

type CompositeGenerator struct {
	generators []Generator
}

func NewCompositeGenerator(generators ...Generator) *CompositeGenerator {
	return &CompositeGenerator{generators: generators}
}

func (g *CompositeGenerator) GenerateForInternalLocation(policies []Policy) GenerateResultFiles {
	var compositeResult GenerateResultFiles

	for _, generator := range g.generators {
		compositeResult = append(compositeResult, generator.GenerateForInternalLocation(policies)...)
	}

	return compositeResult
}

func (g *CompositeGenerator) GenerateForServer(policies []Policy, server http.Server) GenerateResultFiles {
	var compositeResult GenerateResultFiles

	for _, generator := range g.generators {
		compositeResult = append(compositeResult, generator.GenerateForServer(policies, server)...)
	}

	return compositeResult
}

func (g *CompositeGenerator) GenerateForLocation(policies []Policy, location http.Location) GenerateResultFiles {
	var compositeResult GenerateResultFiles

	for _, generator := range g.generators {
		compositeResult = append(compositeResult, generator.GenerateForLocation(policies, location)...)
	}

	return compositeResult
}
