package bridge

import (
	"fmt"
	"log"
	"path/filepath"
)

type ITypeLibrary interface {
	AddType(pkg string, name string, t *Type) (alt *Type)
	GetType(pkg string, name string) *Type
	AddBasicType(name string) (alt *Type)
	GetBasicType(name string) (alt *Type)

	// Package related API
	AddPackage(name string) (shortName string)
	ForEach(func(string, *Type) bool)

	// Signature string creation
	Signature(t *Type) string
	TypeListSignature(types []*Type, argfmt string) string
}

type TypeLibrary struct {
	types       map[string]*Type
	typeCounter int64

	// Package related API
	shortNamesForPkg map[string]string
	pkgByShortName   map[string]string
	pkgCounter       int
}

func NewTypeLibrary() *TypeLibrary {
	log.Println("Creating new package library... ")
	out := TypeLibrary{}
	out.pkgByShortName = make(map[string]string)
	out.shortNamesForPkg = make(map[string]string)

	log.Println("Creating new type library... ")
	out.types = make(map[string]*Type)

	// add some basic types
	out.AddBasicType("error")
	out.AddBasicType("string")
	out.AddBasicType("float")
	out.AddBasicType("float32")
	out.AddBasicType("float64")
	out.AddBasicType("bool")
	out.AddBasicType("byte")
	out.AddBasicType("int")
	out.AddBasicType("int8")
	out.AddBasicType("int16")
	out.AddBasicType("int32")
	out.AddBasicType("int64")
	out.AddBasicType("uint")
	out.AddBasicType("uint8")
	out.AddBasicType("uint16")
	out.AddBasicType("uint32")
	out.AddBasicType("uint64")
	return &out
}

/**
 * Returns all the types as a list.
 */
func (tl *TypeLibrary) ForEach(mapFunc func(string, *Type) bool) {
	for k, v := range tl.types {
		if mapFunc(k, v) {
			return
		}
	}
}

func (tl *TypeLibrary) AddBasicType(name string) (alt *Type) {
	return tl.AddType("", name, &Type{BasicType, name})
}

func (tl *TypeLibrary) GetBasicType(name string) (alt *Type) {
	return tl.GetType("", name)
}

/**
 * Adds a type to the type system.  If the type already exists then
 * the existing one is returned otherwise a new type is added and returned.
 * Also the type's ID will be set.
 */
func (tl *TypeLibrary) AddType(pkg string, name string, t *Type) (alt *Type) {
	tl.AddPackage(pkg)
	key := pkg + "." + name
	if value, ok := tl.types[key]; ok {
		return value
	}
	tl.typeCounter++
	tl.types[key] = t
	return t
}

func (tl *TypeLibrary) GetType(pkg string, name string) *Type {
	key := pkg + "." + name
	t := tl.types[key]
	if t == nil {
		// perhaps we are passing in the short name instead
		key := tl.PackageByShortName(pkg) + "." + name
		t = tl.types[key]
	}
	return t
}

func (tl *TypeLibrary) AddPackage(pkg string) (shortName string) {
	if value, ok := tl.shortNamesForPkg[pkg]; ok {
		return value
	}
	log.Println("Adding Package: ", pkg)

	// create a name and return it (and save it ofcourse)
	name := ""
	tmppkg := pkg
	for tmppkg != "" {
		name := filepath.Base(tmppkg) + name
		if _, ok := tl.shortNamesForPkg[tmppkg]; !ok {
			tl.shortNamesForPkg[tmppkg] = name
			tl.pkgByShortName[name] = tmppkg
			return name
		}
		tmppkg = filepath.Dir(tmppkg)
	}

	// return a random name!
	name = fmt.Sprintf("pkg%d", tl.pkgCounter)
	tl.shortNamesForPkg[pkg] = name
	tl.pkgByShortName[name] = pkg
	return name
}

func (tl *TypeLibrary) PackageByShortName(name string) string {
	if value, ok := tl.pkgByShortName[name]; ok {
		return value
	}
	return ""
}

func (tl *TypeLibrary) GetShortPackageName(pkg string) string {
	if value, ok := tl.shortNamesForPkg[pkg]; ok {
		return value
	}
	return ""
}

func (tl *TypeLibrary) Signature(t *Type) string {
	switch t.TypeClass {
	case NullType:
		return ""
	case UnresolvedType:
		return t.TypeData.(string)
	case BasicType:
		return t.TypeData.(string)
	case AliasType:
		return t.TypeData.(string)
	case ReferenceType:
		return "*" + tl.Signature(t.TypeData.(*ReferenceTypeData).TargetType)
	case RecordType:
		return t.TypeData.(*RecordTypeData).Name
	case TupleType:
		out := "("
		for index, childType := range t.TypeData.(*TupleTypeData).SubTypes {
			if index > 0 {
				out += ", "
			}
			out += tl.Signature(childType)
		}
		return out
	case FunctionType:
		funcTypeData := t.TypeData.(*FunctionTypeData)
		out := "func ("
		out += tl.TypeListSignature(funcTypeData.InputTypes, "") + ")"
		if funcTypeData.OutputTypes != nil {
			out += "(" + tl.TypeListSignature(funcTypeData.OutputTypes, "") + ")"
		}
		if funcTypeData.ExceptionTypes != nil {
			out += " throws (" + tl.TypeListSignature(funcTypeData.ExceptionTypes, "") + ")"
		}
		return out
	case ListType:
		return "[]" + tl.Signature(t.TypeData.(*ListTypeData).TargetType)
	case MapType:
		mapTypeData := t.TypeData.(*MapTypeData)
		return "map[" + tl.Signature(mapTypeData.KeyType) + "]" + tl.Signature(mapTypeData.ValueType)
	}
	return ""
}

func (tl *TypeLibrary) TypeListSignature(types []*Type, argfmt string) string {
	out := ""
	if types != nil {
		for index, inType := range types {
			if index > 0 {
				out += ","
			}
			out += fmt.Sprintf(argfmt, index) + " "
			out += tl.Signature(inType)
		}
	}
	return out
}
