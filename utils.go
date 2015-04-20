package bridge

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"reflect"
)

type ParsedFile struct {
	FullPath string
	FileNode *ast.File
	Package  string
	Imports  map[string]string
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
func (parsedFile *ParsedFile) ProcessNode(typeSystem ITypeSystem) error {
	for _, decl := range parsedFile.FileNode.Decls {
		gendecl, ok := decl.(*ast.GenDecl)
		if ok && len(gendecl.Specs) > 0 {
			typeSpec, ok := gendecl.Specs[0].(*ast.TypeSpec)
			if ok {
				NodeToType(typeSpec, parsedFile.FileNode.Name.Name, typeSystem)
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
func NodeToType(node ast.Node, pkg string, typeSystem ITypeSystem) *Type {
	switch typeExpr := node.(type) {
	case *ast.StarExpr:
		// we have a reference type
		out := &Type{TypeClass: ReferenceType}
		out.TypeData = &ReferenceTypeData{TargetType: NodeToType(typeExpr.X, pkg, typeSystem)}
		return out
	case *ast.FuncType:
		{
			out := &Type{TypeClass: FunctionType}

			// create a function type
			functionType := &FunctionTypeData{}
			out.TypeData = functionType
			for _, param := range typeExpr.Params.List {
				paramType := NodeToType(param.Type, pkg, typeSystem)
				functionType.InputTypes = append(functionType.InputTypes, paramType)
			}
			if typeExpr.Results != nil && typeExpr.Results.List != nil {
				for _, result := range typeExpr.Results.List {
					resultType := NodeToType(result.Type, pkg, typeSystem)
					functionType.OutputTypes = append(functionType.OutputTypes, resultType)
				}
			}
			return out
		}
	case *ast.MapType:
		typeData := &MapTypeData{}
		typeData.KeyType = NodeToType(typeExpr.Key, pkg, typeSystem)
		typeData.ValueType = NodeToType(typeExpr.Value, pkg, typeSystem)
		return &Type{TypeClass: MapType, TypeData: typeData}
	case *ast.ArrayType:
		return &Type{TypeClass: ListType,
			TypeData: &ListTypeData{TargetType: NodeToType(typeExpr.Elt, pkg, typeSystem)}}
	case *ast.Ident:
		t := typeSystem.GetType(pkg, typeExpr.Name)
		if t == nil {
			t = &Type{TypeClass: LazyType, TypeData: typeExpr.Name}
			typeSystem.AddType(pkg, typeExpr.Name, t)
		}
		return t
	case *ast.SelectorExpr:
		pkgName := typeExpr.X.(*ast.Ident).Name
		t := typeSystem.GetType(pkgName, typeExpr.Sel.Name)
		if t == nil {
			t = &Type{TypeClass: LazyType, TypeData: typeExpr.Sel.Name}
			typeSystem.AddType(pkgName, typeExpr.Sel.Name, t)
		}
		return t
	case *ast.StructType:
		{
			recordData := &RecordTypeData{}
			fieldList := typeExpr.Fields.List
			for _, field := range fieldList {
				fieldType := NodeToType(field.Type, pkg, typeSystem)
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
				fieldType := NodeToType(field.Type, pkg, typeSystem)
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
		currT := typeSystem.GetType(pkg, recordData.Name)
		if currT == nil {
			typeSystem.AddType(pkg, recordData.Name, out)
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
				t2 := NodeToType(typeExpr, "", typeSystem)
				out.TypeData = t2.TypeData
			}
		case *ast.StructType:
			{
				t2 := NodeToType(typeExpr, "", typeSystem)
				out.TypeData = t2.TypeData
			}
		}
		return out
	}
	log.Println("Damn - the wrong type: ", node, reflect.TypeOf(node))
	return nil
}
