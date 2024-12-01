package main

import (
	"fmt"
	"strings"

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
			generateFile(gen, f)
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File) {
	// Generate Python module path from proto package
	pythonPath := strings.ReplaceAll(string(file.Desc.Package()), ".", "/")
	filename := fmt.Sprintf("%s/%s_mcp.py", pythonPath, file.Desc.Name())

	// Create new generated file without Go import path
	g := gen.NewGeneratedFile(filename, "")

	// Generate Python imports
	g.P("import os")
	g.P("import json")
	g.P("import logging")
	g.P("from typing import Any, Sequence")
	g.P("from mcp.server import Server")
	g.P("from mcp.types import Resource, Tool, TextContent, ImageContent, EmbeddedResource")
	g.P("from pydantic import AnyUrl")
	g.P()

	// Create server instance
	g.P("# Configure logging")
	g.P("logging.basicConfig(level=logging.INFO)")
	g.P(fmt.Sprintf("logger = logging.getLogger('%s-server')", string(file.Desc.Package())))
	g.P()
	g.P(fmt.Sprintf("app = Server('%s-server')", string(file.Desc.Package())))
	g.P()

	// Generate tool definitions for each service
	for _, service := range file.Services {
		generateServiceTools(g, service)
	}

	// Generate main function
	g.P("async def main():")
	g.P("    from mcp.server.stdio import stdio_server")
	g.P("    async with stdio_server() as (read_stream, write_stream):")
	g.P("        await app.run(")
	g.P("            read_stream,")
	g.P("            write_stream,")
	g.P("            app.create_initialization_options()")
	g.P("        )")
	g.P()
	g.P("if __name__ == '__main__':")
	g.P("    import asyncio")
	g.P("    asyncio.run(main())")
}

func generateServiceTools(g *protogen.GeneratedFile, service *protogen.Service) {
	g.P("@app.list_tools()")
	g.P("async def list_tools() -> list[Tool]:")
	g.P("    return [")

	for _, method := range service.Methods {
		g.P("        Tool(")
		g.P(fmt.Sprintf("            name='%s',", method.Desc.Name()))
		g.P(fmt.Sprintf("            description='%s method from %s service',", method.Desc.Name(), service.Desc.Name()))
		g.P("            inputSchema={")
		g.P("                'type': 'object',")
		g.P("                'properties': {")
		// TODO: Generate input schema from method input message
		g.P("                },")
		g.P("                'required': []")
		g.P("            }")
		g.P("        ),")
	}

	g.P("    ]")
	g.P()

	// Generate call_tool handler
	g.P("@app.call_tool()")
	g.P("async def call_tool(name: str, arguments: Any) -> Sequence[TextContent | ImageContent | EmbeddedResource]:")
	g.P("    if name not in [")
	for _, method := range service.Methods {
		g.P(fmt.Sprintf("        '%s',", method.Desc.Name()))
	}
	g.P("    ]:")
	g.P("        raise ValueError(f'Unknown tool: {name}')")
	g.P()
	g.P("    # TODO: Implement method handlers")
	g.P("    return [")
	g.P("        TextContent(")
	g.P("            type='text',")
	g.P("            text=json.dumps({'result': 'Not implemented'}, indent=2)")
	g.P("        )")
	g.P("    ]")
}
