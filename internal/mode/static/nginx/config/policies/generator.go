package policies

import (
	"maps"

	"github.com/nginxinc/nginx-gateway-fabric/internal/mode/static/nginx/config/http"
)

type Generator interface {
	GenerateForLocation(policies []Policy, location http.Location) GenerateResult
	// TODO: do we need this as a separate method?
	GenerateForInternalLocation(policies []Policy, internalLocation http.Location) GenerateResult
	GenerateForServer(policies []Policy, server http.Server) GenerateResult
}

type GenerateResult struct {
	KeyVals map[string]interface{}
	Files   []File
}

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

func (g *CompositeGenerator) GenerateForInternalLocation(
	policies []Policy,
	internalLocation http.Location,
) GenerateResult {
	var compositeResult GenerateResult

	for _, generator := range g.generators {
		result := generator.GenerateForInternalLocation(policies, internalLocation)
		compositeResult.Files = append(compositeResult.Files, result.Files...)
		maps.Copy(compositeResult.KeyVals, result.KeyVals)
	}

	return compositeResult
}

func (g *CompositeGenerator) GenerateForServer(policies []Policy, server http.Server) GenerateResult {
	var compositeResult GenerateResult

	for _, generator := range g.generators {
		result := generator.GenerateForServer(policies, server)
		compositeResult.Files = append(compositeResult.Files, result.Files...)
		maps.Copy(compositeResult.KeyVals, result.KeyVals)
	}

	return compositeResult
}

func (g *CompositeGenerator) GenerateForLocation(policies []Policy, location http.Location) GenerateResult {
	var compositeResult GenerateResult

	for _, generator := range g.generators {
		result := generator.GenerateForLocation(policies, location)
		compositeResult.Files = append(compositeResult.Files, result.Files...)
		maps.Copy(compositeResult.KeyVals, result.KeyVals)
	}

	return compositeResult
}
