package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/panyam/bridge"
	"github.com/panyam/bridge/rest"
	"io"
	"log"
	"os"
)

func main() {
	var serviceName, operation string
	flag.StringVar(&serviceName, "service", "", "The service whose methods are to be extracted and for whome binding code is to be generated")
	flag.StringVar(&operation, "operation", "", "The operation within the service to be generated code for.  If this is empty or not provided then ALL operations in the service will code generated for them")

	flag.Parse()

	if serviceName == "" {
		log.Println("Service required")
	}

	_, typeLibrary := ParseFiles(flag.Args())
	serviceType := typeLibrary.GetType("core", serviceName)
	CreateClientForType(typeLibrary, serviceType)
}

func ParseFiles(fileNames []string) (map[string]*bridge.ParsedFile, bridge.ITypeLibrary) {
	typeLibrary := bridge.NewTypeLibrary()
	parsedFiles := make(map[string]*bridge.ParsedFile)
	for _, srcFile := range fileNames {
		pf, err := bridge.NewParsedFile(srcFile)
		if err != nil {
			log.Println("Parsing error: ", err)
			panic(err)
		}
		parsedFiles[pf.FullPath] = pf
	}

	// parse all files now that we have the imports kind of resolved
	// At this point in each of the file, the Import map should have a local
	// name as well as a path ("." imports are ignored for now)
	for path, parsedFile := range parsedFiles {
		log.Println("Processing Src: ", path)
		parsedFile.ProcessNode(typeLibrary)
	}

	// Report unresolved types
	typeLibrary.ForEach(func(name string, t *bridge.Type, stop *bool) {
		if t.TypeClass == bridge.UnresolvedType {
			log.Println("Unresolved Type: ", name, t.TypeData)
		}
	})
	return parsedFiles, typeLibrary
}

func OpenFile(path string) *os.File {
	out, err := os.Create("./restclient/client.go")
	if err != nil {
		log.Println("Cannot create file: ", err)
		panic(err)
	}
	return out
}

func CreateClientForType(typeLibrary bridge.ITypeLibrary, serviceType *bridge.Type) {
	// Create the generator
	generator := rest.NewGenerator(nil, typeLibrary, "../rest/templates/")

	sigVisited := make(map[string]bool)
	typeVisited := make(map[*bridge.Type]bool)
	uniqueTypes := make([]*bridge.Type, 0, 100)
	resetTypes := func() {
		sigVisited = make(map[string]bool)
		typeVisited = make(map[*bridge.Type]bool)
		uniqueTypes = make([]*bridge.Type, 0, 100)
	}
	generator.TypeMarker = func(t *bridge.Type) {
		if !typeVisited[t] {
			typeVisited[t] = true
			sig := typeLibrary.Signature(t)
			if !sigVisited[sig] {
				sigVisited[sig] = true
				uniqueTypes = append(uniqueTypes, t)
			}
		}
	}

	// Generate the interface declartion
	resetTypes()
	clientBuff := bytes.NewBuffer(nil)
	err := generator.EmitClientClass(clientBuff, serviceType)
	if err != nil {
		log.Println("Class emitting error: ", err)
		return
	}
	client_file := OpenFile("./restclient/client.go")
	EmitFileHeader(client_file, generator.ClientPackageName, uniqueTypes, typeLibrary)
	client_file.Write(clientBuff.Bytes())
	client_file.Close()

	// Generate code for each of the service operation methods
	resetTypes()
	opsBuff := bytes.NewBuffer(nil)
	serviceTypeData := generator.ServiceTypeData()
	for _, field := range serviceTypeData.Fields {
		switch optype := field.Type.TypeData.(type) {
		case *bridge.FunctionTypeData:
			// get the type info and ensure the packages referred by this type
			// are imported
			generator.EmitSendRequestMethod(opsBuff, field.Name, optype, "arg")
		}
	}
	ops_file := OpenFile("./restops/ops.go")
	EmitFileHeader(ops_file, generator.ClientPackageName, uniqueTypes, typeLibrary)
	ops_file.Write(opsBuff.Bytes())
	ops_file.Close()

	// Write the writers for each of the unique types
	resetTypes()
	writersBuff := bytes.NewBuffer(nil)
	for _, t := range uniqueTypes {
		generator.EmitTypeWriter(writersBuff, t)
	}
	writers_file := OpenFile("./restwriters/writers.go")
	EmitFileHeader(writers_file, generator.ClientPackageName, uniqueTypes, typeLibrary)
	writers_file.Write(writersBuff.Bytes())
	writers_file.Close()
}

/**
 * Writes the package header containing the package name and the imports of the
 * unique types to the output.
 */
func EmitFileHeader(writer io.Writer, packageName string, types []*bridge.Type, typeLib bridge.ITypeLibrary) error {
	writer.Write([]byte("package " + packageName + "\n\n"))

	writer.Write([]byte("import (\n"))
	pkgVisited := make(map[string]bool)
	for _, t := range types {
		leafType := t.LeafType()
		pkg := leafType.Package
		if leafType != nil && pkg != "" && !pkgVisited[pkg] {
			pkgVisited[pkg] = true
			writer.Write([]byte(fmt.Sprintf("	%s \"%s\"\n", typeLib.ShortNameForPackage(pkg), pkg)))
		}
	}
	writer.Write([]byte(")\n"))
	return nil
}
