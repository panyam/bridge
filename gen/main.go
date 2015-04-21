package main

import (
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

	/*
		var parsedFile *ast.File = nil
		fset := token.NewFileSet() // positions are relative to fset
		newPkg, err := ast.NewPackage(fset, parsedFiles, nil, nil)
		log.Println("NP: ", newPkg, err)
		parsedFile = ast.MergePackageFiles(newPkg, 0)
		log.Println("Big ParsedFile: ", parsedFile)
		log.Println("Package: ", parsedFile.Name.Name)
		log.Println("Imports: ", parsedFile.Imports)
		log.Println("UnResolved: ", parsedFile.Unresolved)
	*/

	typeLibrary.ForEach(func(name string, t *bridge.Type) bool {
		if t.TypeClass == bridge.UnresolvedType {
			log.Println("Unresolved Type: ", name, t.TypeData)
		}
		return false
	})

	generator := rest.NewGenerator(nil, typeLibrary, "../rest/templates/")
	serviceType := typeLibrary.GetType("core", serviceName)
	err := generator.EmitClientClass(serviceType)
	if err == nil {
		serviceTypeData := generator.ServiceTypeData()
		// Generate code for each of the service methods
		for _, field := range serviceTypeData.Fields {
			switch optype := field.Type.TypeData.(type) {
			case *bridge.FunctionTypeData:
				generator.EmitSendRequestMethod(os.Stdout, field.Name, optype, "arg")
			}
		}
	} else {
		log.Println("Class emitting error: ", err)
	}
}
