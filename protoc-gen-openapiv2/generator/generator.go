package generator

import (
	"github.com/longXboy/grpc-gateway/v2/internal/codegenerator"
	"github.com/longXboy/grpc-gateway/v2/internal/descriptor"
	"github.com/longXboy/grpc-gateway/v2/protoc-gen-openapiv2/internal/genopenapi"
	"google.golang.org/protobuf/types/pluginpb"
)

// Generator is openapi v2 generator
type Generator struct {
}

// Gen generates openapi v2 json content
func (g *Generator) Gen(req *pluginpb.CodeGeneratorRequest, onlyRPC bool) (*pluginpb.CodeGeneratorResponse, error) {
	reg := descriptor.NewRegistry()
	reg.SetUseJSONNamesForFields(true)
	reg.SetRecursiveDepth(1024)
	reg.SetMergeFileName("apidocs")
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
