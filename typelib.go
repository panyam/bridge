package bridge

import (
	"fmt"
	"log"
	"path/filepath"
)

type ITypeLibrary interface {
	AddType(pkg string, name string, t *Type) (alt *Type)
	GetType(pkg string, name string) *Type
	GetShortPackageName(pkg string) string
}

type TypeLibrary struct {
	shortNamesForPkg map[string]string
	pkgByShortName   map[string]string
	types            map[string]*Type
	typeCounter      int64
	pkgCounter       int
}

func NewTypeLibrary() *TypeLibrary {
	log.Println("Creating new type library... ")
	out := TypeLibrary{}
	out.pkgByShortName = make(map[string]string)
	out.shortNamesForPkg = make(map[string]string)
	out.types = make(map[string]*Type)

	// add some basic types
	out.AddBasicType("string")
	out.AddBasicType("float")
	out.AddBasicType("float32")
	out.AddBasicType("float64")
	out.AddBasicType("bool")
	out.AddBasicType("byte")
	out.AddBasicType("int")
	out.AddBasicType("int16")
	out.AddBasicType("int32")
	out.AddBasicType("int64")
	out.AddBasicType("uint")
	out.AddBasicType("uint16")
	out.AddBasicType("uint32")
	out.AddBasicType("uint64")
	return &out
}

func (tl *TypeLibrary) AddBasicType(name string) (alt *Type) {
	return tl.AddType("", name, &Type{BasicType, name})
}

/**
 * Adds a type to the type system.  If the type already exists then
 * the existing one is returned otherwise a new type is added and returned.
 * Also the type's ID will be set.
 */
func (tl *TypeLibrary) AddType(pkg string, name string, t *Type) (alt *Type) {
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
	return tl.types[key]
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

	// create a name and return it (and save it ofcourse)
	name := ""
	for pkg != "" {
		name := filepath.Base(pkg) + name
		if _, ok := tl.shortNamesForPkg[pkg]; !ok {
			tl.shortNamesForPkg[name] = pkg
			tl.pkgByShortName[pkg] = name
			return name
		}
		pkg = filepath.Dir(pkg)
	}

	// return a random name!
	name = fmt.Sprintf("pkg%d", tl.pkgCounter)
	tl.shortNamesForPkg[name] = pkg
	tl.pkgByShortName[pkg] = name
	return name
}
