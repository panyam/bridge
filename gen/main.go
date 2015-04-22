package main

import (
	"bytes"
	"flag"
	"github.com/panyam/bridge"
	"github.com/panyam/bridge/rest"
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

	typeLibrary := bridge.NewTypeLibrary()
	parsedFiles := make(map[string]*bridge.ParsedFile)
	for _, srcFile := range flag.Args() {
		pf, err := bridge.NewParsedFile(srcFile)
		if err != nil {
			log.Println("Parsing error: ", err)
			return
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

	serviceType := typeLibrary.GetType("core", serviceName)
	CreateClientForType(typeLibrary, serviceType)
}

func CreateClientForType(typeLibrary bridge.ITypeLibrary, serviceType *bridge.Type) {
	client_file, err := os.Create("./restclient/client.go")
	if err != nil {
		log.Println("Cannot create client_file: ", err)
		return
	}
	defer client_file.Close()

	opsbuffer := bytes.NewBuffer(nil)
	generator := rest.NewGenerator(nil, typeLibrary, "../rest/templates/")
	err = generator.EmitClientClass(client_file, serviceType)
	if err != nil {
		log.Println("Class emitting error: ", err)
		return
	}

	uniqueTypes := make(map[*bridge.Type]bool)

	serviceTypeData := generator.ServiceTypeData()
	// Generate code for each of the service methods
	for _, field := range serviceTypeData.Fields {
		switch optype := field.Type.TypeData.(type) {
		case *bridge.FunctionTypeData:
			// get the type info and ensure the packages referred by this type
			// are imported
			generator.EmitSendRequestMethod(opsbuffer, field.Name, optype, "arg")
			for _, t := range optype.InputTypes {
				uniqueTypes[t] = true
			}
			for _, t := range optype.OutputTypes {
				uniqueTypes[t] = true
			}
		}
	}

	operations_file, err := os.Create("./restclient/ops.go")
	if err != nil {
		log.Println("Cannot create operations_file: ", err)
		return
	}
	defer operations_file.Close()
	operations_file.Write(opsbuffer.Bytes())

	// Now generate the writers for all types we have found
	writers_file, err := os.Create("./restclient/writers.go")
	if err != nil {
		log.Println("Cannot create writers_file: ", err)
		return
	}
	defer writers_file.Close()
	visited := make(map[*bridge.Type]bool)
	for t, _ := range uniqueTypes {
		log.Println("Creating writer for T: ", t)
		generator.EmitTypeWriter(writers_file, t, visited)
	}

	/*
		headerbuffer := bytes.NewBuffer(nil)
		headerbuffer.Write([]byte(fmt.Sprintf("package %s\n", generator.ClientPackageName)))
		headerbuffer.Write([]byte(fmt.Sprintf("import (\n")
		headerbuffer.Write([]byte(fmt.Sprintf(")\n\n")
	*/
}
