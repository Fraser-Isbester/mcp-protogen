package main

import (
	"github.com/fraser-isbester/protoc-gen-mcp/internal/generator"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			if err := generator.Generate(gen, f); err != nil {
				return err
			}
		}
		return nil
	})
}
