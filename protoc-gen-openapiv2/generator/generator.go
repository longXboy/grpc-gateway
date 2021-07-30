package generator

import (
	"github.com/longXboy/grpc-gateway/v2/internal/codegenerator"
	"github.com/longXboy/grpc-gateway/v2/internal/descriptor"
	"github.com/longXboy/grpc-gateway/v2/protoc-gen-openapiv2/internal/genopenapi"
	"google.golang.org/protobuf/types/pluginpb"
)

// Generator is openapi v2 generator
type Generator struct {
	reg *descriptor.Registry
}

type Option func(gen *Generator)

// UseJSONNamesForFields. if disabled, the original proto name will be used for generating OpenAPI definitions
func UseJSONNamesForFields(b bool) Option {
	return func(gen *Generator) {
		gen.reg.SetUseJSONNamesForFields(b)
	}
}

// RecursiveDepth. maximum recursion count allowed for a field type
func RecursiveDepth(depth int) Option {
	return func(gen *Generator) {
		gen.reg.SetRecursiveDepth(depth)
	}
}

// EnumsAsInts. whether to render enum values as integers, as opposed to string values
func EnumsAsInts(b bool) Option {
	return func(gen *Generator) {
		gen.reg.SetEnumsAsInts(b)
	}
}

// MergeFileName. target OpenAPI file name prefix after merge
func MergeFileName(name string) Option {
	return func(gen *Generator) {
		gen.reg.SetMergeFileName(name)
	}
}

// DisableDefaultErrors. if set, disables generation of default errors. This is useful if you have defined custom error handling
func DisableDefaultErrors(b bool) Option {
	return func(gen *Generator) {
		gen.reg.SetDisableDefaultErrors(b)
	}
}

func NewGenerator(options ...Option) *Generator {
	gen := &Generator{
		reg: descriptor.NewRegistry(),
	}
	gen.reg.SetUseJSONNamesForFields(true)
	gen.reg.SetRecursiveDepth(1024)
	gen.reg.SetMergeFileName("apidocs")
	gen.reg.SetDisableDefaultErrors(true)
	for _, o := range options {
		o(gen)
	}
	return gen
}

// Gen generates openapi v2 json content
func (g *Generator) Gen(req *pluginpb.CodeGeneratorRequest, onlyRPC bool) (*pluginpb.CodeGeneratorResponse, error) {
	reg := g.reg
	if reg == nil {
		reg = NewGenerator().reg
	}
	reg.SetGenerateRPCMethods(onlyRPC)
	if err := reg.SetRepeatedPathParamSeparator("csv"); err != nil {
		return nil, err
	}

	gen := genopenapi.New(reg)

	if err := genopenapi.AddErrorDefs(reg); err != nil {
		return nil, err
	}

	if err := reg.Load(req); err != nil {
		return nil, err
	}
	var targets []*descriptor.File
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			return nil, err
		}
		targets = append(targets, f)
	}
	out, err := gen.Generate(targets)
	if err != nil {
		return nil, err
	}
	return emitFiles(out), nil
}

func emitFiles(out []*descriptor.ResponseFile) *pluginpb.CodeGeneratorResponse {
	files := make([]*pluginpb.CodeGeneratorResponse_File, len(out))
	for idx, item := range out {
		files[idx] = item.CodeGeneratorResponse_File
	}
	resp := &pluginpb.CodeGeneratorResponse{File: files}
	codegenerator.SetSupportedFeaturesOnCodeGeneratorResponse(resp)
	return resp
}
