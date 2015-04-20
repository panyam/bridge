package bridge

import (
	"log"
)

type ITypeLibrary interface {
	AddType(pkg string, name string, t *Type) (alt *Type)
	GetType(pkg string, name string) *Type
}

type TypeLibrary struct {
	types       map[string]*Type
	typeCounter int64
}

func NewTypeLibrary() *TypeLibrary {
	log.Println("Creating new type system... ", pkg)
	out := TypeLibrary{}
	out.types = make(map[string]*Type)
	return &out
}

/**
 * Adds a type to the type system.  If the type already exists then
 * the existing one is returned otherwise a new type is added and returned.
 * Also the type's ID will be set.
 */
func (ts *TypeLibrary) AddType(pkg string, name string, t *Type) (alt *Type) {
	key := pkg + "." + name
	if value, ok := ts.types[key]; ok {
		return value
	}
	ts.typeCounter++
	ts.types[key] = t
	return t
}

func (ts *TypeLibrary) GetType(pkg string, name string) *Type {
	key := pkg + "." + name
	return ts.types[key]
}

func (ts *TypeLibrary) FindType(typeClass int, typeData interface{}) string {
	return ""
}
