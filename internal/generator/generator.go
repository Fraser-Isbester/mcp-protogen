package generator

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

type Field struct {
	Name        string
	Type        string
	Description string
}

type Method struct {
	Name           string
	Description    string
	InputFields    []Field
	RequiredFields []string
}

type Service struct {
	Name    string
	Methods []Method
}

type TemplateData struct {
	Package  string
	Services []Service
}

// ConvertType converts protobuf types to JSON schema types
func ConvertType(field protoreflect.FieldDescriptor) string {
	switch field.Kind() {
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BoolKind:
		return "boolean"
	case protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind:
		return "number"
	default:
		return "string"
	}
}

func Generate(gen *protogen.Plugin, file *protogen.File) error {
	data := TemplateData{
		Package:  string(file.Desc.Package()),
		Services: make([]Service, 0, len(file.Services)),
	}

	for _, service := range file.Services {
		s := Service{
			Name:    string(service.Desc.Name()),
			Methods: make([]Method, 0, len(service.Methods)),
		}

		for _, method := range service.Methods {
			m := Method{
				Name:           string(method.Desc.Name()),
				Description:    fmt.Sprintf("%s method from %s service", method.Desc.Name(), service.Desc.Name()),
				InputFields:    make([]Field, 0),
				RequiredFields: make([]string, 0),
			}

			// Process input message fields
			fields := method.Input.Desc.Fields()
			for i := 0; i < fields.Len(); i++ {
				field := fields.Get(i)
				f := Field{
					Name:        string(field.Name()),
					Type:        ConvertType(field),
					Description: fmt.Sprintf("Field %s", field.Name()),
				}
				m.InputFields = append(m.InputFields, f)

				if !field.HasOptionalKeyword() {
					m.RequiredFields = append(m.RequiredFields, string(field.Name()))
				}
			}

			s.Methods = append(s.Methods, m)
		}

		data.Services = append(data.Services, s)
	}

	// Load and parse template
	tmplContent, err := templateFS.ReadFile("templates/mcp.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read template file: %v", err)
	}

	tmpl, err := template.New("mcp").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Generate the file with the correct name
	packageParts := strings.Split(string(file.Desc.Package()), ".")
	baseName := packageParts[0]

	filename := fmt.Sprintf("%s/%s_mcp.py",
		strings.ReplaceAll(string(file.Desc.Package()), ".", "/"),
		baseName)
	g := gen.NewGeneratedFile(filename, "")
	_, err = g.Write(buf.Bytes())
	return err
}
