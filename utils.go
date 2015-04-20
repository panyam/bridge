package bridge

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

type ParsedFile struct {
	FullPath    string
	FileNode    *ast.File
	Package     string
	PackagePath string
	Imports     map[string]string
}

// Given a full path to a file or folder, finds the full name of the
// package path in one of the GOPATH or GOROOT folders
// For example, if GOPATH contained the following folders:
// /home/user1/a
// /home/user1/b
//
// and GOROOT was:
// /user/local/go
//
// and given the input full path of:
//
// /home/user1/a/src/github.com/theuser/repo
//
// would return
//
// github.com/theuser/repo
//
// as the full package path as it exists in the src folder of one of the folders
// in GOPATH or GOROOT
func PackagePathForFile(fullpath string) string {
	gopath := strings.Split(os.Getenv("GOPATH"), ";")
	folders := []string{runtime.GOROOT()}
	folders = append(folders, gopath...)
	for _, folder := range folders {
		if strings.HasPrefix(fullpath, folder+"src/") {
			return fullpath[len(folder+"src/"):]
		}
	}
	return ""
}

func NewParsedFile(srcFile string) (out *ParsedFile, err error) {
	out = &ParsedFile{Imports: make(map[string]string)}
	srcFile, err = filepath.Abs(srcFile)
	if err != nil {
		return nil, err
	}
	out.FullPath = srcFile
	fset := token.NewFileSet() // positions are relative to fset
	out.FileNode, err = parser.ParseFile(fset, srcFile, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, err
	}
	out.Package = out.FileNode.Name.Name
	out.PackagePath = PackagePathForFile(filepath.Dir(srcFile))
	out.Imports[out.Package] = out.PackagePath
	for _, importSpec := range out.FileNode.Imports {
		path := importSpec.Path.Value
		name := ""
		if importSpec.Name == nil {
			name = filepath.Base(path)
		} else {
			name = importSpec.Name.Name
		}

		if name == "." {
			log.Println("Not sure how to deal with . imports", path)
		} else {
			out.Imports[name] = path
		}
	}
	return out, err
}

/**
 * Parses a file and returns a map of types indexed by name.
 */
func (parsedFile *ParsedFile) ProcessNode(typeLibrary ITypeLibrary) error {
	for _, decl := range parsedFile.FileNode.Decls {
		gendecl, ok := decl.(*ast.GenDecl)
		if ok && len(gendecl.Specs) > 0 {
			typeSpec, ok := gendecl.Specs[0].(*ast.TypeSpec)
			if ok {
				parsedFile.NodeToType(typeSpec, typeLibrary)
			}
		}
	}
	return nil
}

/**
 * Finds a GenDecl node in a parsed file.
 */
func FindDecl(parsedFile *ast.File, declName string) *ast.GenDecl {
	for _, decl := range parsedFile.Decls {
		gendecl, ok := decl.(*ast.GenDecl)
		if ok && len(gendecl.Specs) > 0 {
			typeSpec := gendecl.Specs[0].(*ast.TypeSpec)
			if ok && typeSpec.Name.Name == declName {
				return gendecl
			}
		}
	}
	return nil
}

/**
 * Convert a node to a type.
 */
func (parsedFile *ParsedFile) NodeToType(node ast.Node, typeLibrary ITypeLibrary) *Type {
	switch typeExpr := node.(type) {
	case *ast.StarExpr:
		// we have a reference type
		out := &Type{TypeClass: ReferenceType}
		out.TypeData = &ReferenceTypeData{TargetType: parsedFile.NodeToType(typeExpr.X, typeLibrary)}
		return out
	case *ast.FuncType:
		{
			out := &Type{TypeClass: FunctionType}

			// create a function type
			functionType := &FunctionTypeData{}
			out.TypeData = functionType
			for _, param := range typeExpr.Params.List {
				paramType := parsedFile.NodeToType(param.Type, typeLibrary)
				functionType.InputTypes = append(functionType.InputTypes, paramType)
			}
			if typeExpr.Results != nil && typeExpr.Results.List != nil {
				for _, result := range typeExpr.Results.List {
					resultType := parsedFile.NodeToType(result.Type, typeLibrary)
					functionType.OutputTypes = append(functionType.OutputTypes, resultType)
				}
			}
			return out
		}
	case *ast.MapType:
		typeData := &MapTypeData{}
		typeData.KeyType = parsedFile.NodeToType(typeExpr.Key, typeLibrary)
		typeData.ValueType = parsedFile.NodeToType(typeExpr.Value, typeLibrary)
		return &Type{TypeClass: MapType, TypeData: typeData}
	case *ast.ArrayType:
		return &Type{TypeClass: ListType,
			TypeData: &ListTypeData{TargetType: parsedFile.NodeToType(typeExpr.Elt, typeLibrary)}}
	case *ast.Ident:
		t := typeLibrary.GetType(parsedFile.PackagePath, typeExpr.Name)
		if t == nil {
			t = &Type{TypeClass: LazyType, TypeData: typeExpr.Name}
			typeLibrary.AddType(parsedFile.PackagePath, typeExpr.Name, t)
		}
		return t
	case *ast.SelectorExpr:
		pkgName := typeExpr.X.(*ast.Ident).Name
		t := typeLibrary.GetType(parsedFile.Imports[pkgName], typeExpr.Sel.Name)
		if t == nil {
			t = &Type{TypeClass: LazyType, TypeData: typeExpr.Sel.Name}
			typeLibrary.AddType(parsedFile.Imports[pkgName], typeExpr.Sel.Name, t)
		}
		return t
	case *ast.StructType:
		{
			recordData := &RecordTypeData{}
			fieldList := typeExpr.Fields.List
			for _, field := range fieldList {
				fieldType := parsedFile.NodeToType(field.Type, typeLibrary)
				// log.Println("Processing field: ", index, field.Names, field.Type, reflect.TypeOf(field.Type))
				if len(field.Names) == 0 {
					recordData.Bases = append(recordData.Bases, fieldType)
				} else {
					for _, fieldName := range field.Names {
						field := &Field{Name: fieldName.Name, Type: fieldType}
						recordData.Fields = append(recordData.Fields, field)
					}
				}
			}
			return &Type{TypeClass: RecordType, TypeData: recordData}
		}
	case *ast.InterfaceType:
		{
			recordData := &RecordTypeData{}
			fieldList := typeExpr.Methods.List
			for _, field := range fieldList {
				// log.Println("Processing method: ", index, field.Names[0], field.Type, reflect.TypeOf(field.Type))
				fieldType := parsedFile.NodeToType(field.Type, typeLibrary)
				for _, fieldName := range field.Names {
					field := &Field{Name: fieldName.Name, Type: fieldType}
					recordData.Fields = append(recordData.Fields, field)
				}
			}
			return &Type{TypeClass: RecordType, TypeData: recordData}
		}
	case *ast.TypeSpec:
		out := &Type{}
		recordData := &RecordTypeData{Name: typeExpr.Name.Name}
		currT := typeLibrary.GetType(parsedFile.PackagePath, recordData.Name)
		if currT == nil {
			typeLibrary.AddType(parsedFile.PackagePath, recordData.Name, out)
		} else {
			out = currT
			if currT.TypeClass != LazyType {
				// what if it already exists?
				log.Println("ERROR: Redefinition of type: ", recordData.Name, currT.TypeClass)
			}
		}
		out.TypeClass = RecordType
		out.TypeData = recordData

		switch typeExpr := typeExpr.Type.(type) {
		case *ast.InterfaceType:
			{
				t2 := parsedFile.NodeToType(typeExpr, typeLibrary)
				out.TypeData = t2.TypeData
			}
		case *ast.StructType:
			{
				t2 := parsedFile.NodeToType(typeExpr, typeLibrary)
				out.TypeData = t2.TypeData
			}
		}
		return out
	}
	log.Println("Damn - the wrong type: ", node, reflect.TypeOf(node))
	return nil
}
