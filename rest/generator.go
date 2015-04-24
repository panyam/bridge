package rest

import (
	"errors"
	"fmt"
	"github.com/panyam/bridge"
	"io"
	"text/template"
)

/**
 * Responsible for generating the code for the client classes.
 */
type Generator struct {
	// where the templates are
	Bindings     map[string]*HttpBinding
	TypeLib      bridge.ITypeLibrary
	TemplatesDir string

	// Parameters to determine Generated output
	Package           string
	ClientPackageName string
	ServiceName       string
	ClientPrefix      string
	ClientSuffix      string
	httpBindings      map[string]*HttpBinding
	ServiceType       *bridge.Type
	TransportRequest  string
	OpName            string
	OpType            *bridge.FunctionTypeData
	OpMethod          string
	OpEndpoint        string

	// Callbacks from the template to mark certain items in the code generation
	TypeMarker func(types ...*bridge.Type)
}

func (g *Generator) MarkTypes(types []*bridge.Type) string {
	if types != nil {
		return g.MarkType(types...)
	}
	return ""
}

func (g *Generator) MarkType(types ...*bridge.Type) string {
	if g.TypeMarker != nil {
		g.TypeMarker(types...)
	}
	return ""
}

func (g *Generator) ClientName() string {
	return g.ClientPrefix + g.ServiceName + g.ClientSuffix
}

func (g *Generator) ServiceTypeData() *bridge.RecordTypeData {
	return g.ServiceType.TypeData.(*bridge.RecordTypeData)
}

func NewGenerator(bindings map[string]*HttpBinding, typeLib bridge.ITypeLibrary, templatesDir string) *Generator {
	if bindings == nil {
		bindings = make(map[string]*HttpBinding)
	}
	out := Generator{Bindings: bindings,
		TypeLib:           typeLib,
		TemplatesDir:      templatesDir,
		ClientPackageName: "restclient",
		ClientSuffix:      "Client",
		TransportRequest:  "*http.Request",
	}
	return &out
}

/**
 * Emits the class that acts as a client for the service.
 */
func (g *Generator) EmitClientClass(writer io.Writer, serviceType *bridge.Type) error {
	serviceTypeData, ok := serviceType.TypeData.(*bridge.RecordTypeData)
	if !ok {
		return errors.New("Can only classes for record/container types")
	}
	g.ServiceType = serviceType
	g.ServiceName = serviceTypeData.Name

	tmpl, err := template.New("client.gen").ParseFiles(g.TemplatesDir + "client.gen")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(writer, g)
	if err != nil {
		panic(err)
	}
	return err
}

/**
 * For a given service operation, emits a method which:
 * 1. Has inputs the same as those of the underlying service operation,
 * 2. creates a transport level request
 * 3. Sends the transport level request
 * 4. Gets a response from the transport level and returns it
 */
func (g *Generator) EmitSendRequestMethod(output io.Writer, opName string, opType *bridge.FunctionTypeData, argPrefix string) error {
	g.OpName = opName
	g.OpType = opType
	g.OpMethod = "GET"
	g.OpEndpoint = "http://hello.world/"
	g.StartWritingMethod(output, opName, opType, "arg")
	if opType.NumInputs() > 0 {
		if opType.NumInputs() == 1 {
			g.EmitObjectWriterCall(output, nil, "arg0", opType.InputTypes[0])
		} else {
			g.StartWritingList(output)
			for index, param := range opType.InputTypes {
				g.EmitObjectWriterCall(output, index, fmt.Sprintf("arg%d", index), param)
			}
			g.EndWritingList(output)
		}
	}
	g.EndWritingMethod(output, opName, opType)
	return nil
}

func (g *Generator) StartWritingMethod(output io.Writer, opName string, opType *bridge.FunctionTypeData, argPrefix string) error {
	return bridge.RenderTemplate(output, g.TemplatesDir+"/sender_header.gen", g)
}

func (g *Generator) EndWritingMethod(output io.Writer, opName string, opType *bridge.FunctionTypeData) error {
	return bridge.RenderTemplate(output, g.TemplatesDir+"/sender_footer.gen", g)
}

func (g *Generator) IOMethodForType(t *bridge.Type) string {
	switch typeData := t.TypeData.(type) {
	case string:
		return typeData
	case *bridge.NamedTypeData:
		if typeData.Package == "" {
			return typeData.Name
		} else {
			return g.TypeLib.ShortNameForPackage(typeData.Package) + "_" + typeData.Name
		}
	case *bridge.AliasTypeData:
		return g.IOMethodForType(typeData.TargetType)
	case *bridge.ReferenceTypeData:
		return "Ref_" + g.IOMethodForType(typeData.TargetType)
	case *bridge.FunctionTypeData:
		panic(errors.New("Function types cannot be serialized"))
	case *bridge.TupleTypeData:
		panic(errors.New("Warning: Tuple types not supported in GO"))
		return "Tuple"
	case *bridge.RecordTypeData:
		if typeData.Name == "" {
			return "interface"
		}
		if typeData.Package == "" {
			return typeData.Name
		} else {
			return g.TypeLib.ShortNameForPackage(typeData.Package) + "_" + typeData.Name
		}
	case *bridge.MapTypeData:
		return "Map_" + g.IOMethodForType(typeData.KeyType) + "_" + g.IOMethodForType(typeData.ValueType)
	case *bridge.ListTypeData:
		return "List_" + g.IOMethodForType(typeData.TargetType)
	}
	return fmt.Sprintf("UnknownWriter, Type: %d", t.TypeClass)
}

/**
 * Emits the code required to invoke the serializer of an object of a given
 * type.
 */
func (g *Generator) EmitObjectWriterCall(output io.Writer, key interface{}, argName string, argType *bridge.Type) error {
	callString := g.IOMethodForType(argType)
	output.Write([]byte("Write_" + callString + "(body, " + argName + ")\n"))
	return nil
}

/**
 * Emits the code required to start a list.
 */
func (g *Generator) StartWritingList(output io.Writer) {
	output.Write([]byte("body.Write([]byte(\"[\"))\n"))
}

/**
 * Emits the code required to end a list.
 */
func (g *Generator) EndWritingList(output io.Writer) {
	output.Write([]byte("body.Write([]byte(\"]\"))\n"))
}

/**
 * Emits the writer for a particular type and in the process returns via the
 * recorder the types that for which writers must or will be defined.
 */
func (g *Generator) EmitTypeWriter(writer io.Writer, argType *bridge.Type) error {
	tmplType := ""
	switch argType.TypeClass {
	case bridge.ListType:
		tmplType = "list"
	case bridge.MapType:
		tmplType = "map"
	case bridge.ReferenceType:
		tmplType = "ref"
	case bridge.RecordType:
		tmplType = "record"
	case bridge.AliasType:
		tmplType = "alias"
	}
	tmplPath := fmt.Sprintf("%s/writer_%s.gen", g.TemplatesDir, tmplType)
	context := map[string]interface{}{"Gen": g, "Type": argType}
	bridge.RenderTemplate(writer, g.TemplatesDir+"/writer_header.gen", context)
	bridge.RenderTemplate(writer, tmplPath, context)
	return bridge.RenderTemplate(writer, g.TemplatesDir+"/writer_footer.gen", context)
}

/**
 * For a given service operation, emits a method:
 * 1. whose input is a http.Response object
 * 2. Which can be parsed into the output values as expected by the service
 * 	  operations's output signature
 */
/*
func (g *Generator) EmitReadResponseMethod(opName string, opType *bridge.FunctionTypeData, argPrefix string) error {
	g.StartReadingMethod(opName, opType, "arg")
	if opType.NumOutputs() > 0 {
		if opType.NumOutputs() == 1 {
			g.EmitObjectReaderCall("arg0", opType.OutputTypes[0])
		} else {
			g.StartReadingList()
			for index, param := range opType.OutputTypes {
				g.StartReadingChild()
				g.EmitObjectReaderCall(fmt.Sprintf("arg%d", index), param)
				g.EndReadingChild()
			}
			g.EndReadingList()
		}
	}
	g.EndReadingMethod(opName, opType)
}
*/
