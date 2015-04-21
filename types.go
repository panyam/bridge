package bridge

import ()

const (
	NullType = iota
	UnresolvedType
	BasicType
	ExternalType
	AliasType
	ReferenceType
	TupleType
	RecordType
	FunctionType
	ListType
	MapType
)

type Type struct {
	// One of the above type classes
	TypeClass int

	TypeData interface{}
}

func NewType(typeCls int, data interface{}) *Type {
	return &Type{TypeClass: typeCls, TypeData: data}
}

type ExternalTypeData struct {
	Package string
	Name    string
}

type AliasTypeData struct {
	// Type this is an alias/typedef for
	Name     string
	AliasFor *Type
}

type ReferenceTypeData struct {
	// The target type this is a reference to
	TargetType *Type
}

type MapTypeData struct {
	// The target type this is an array of
	KeyType   *Type
	ValueType *Type
}

type ListTypeData struct {
	// The target type this is an array of
	TargetType *Type
}

type TupleTypeData struct {
	SubTypes []*Type
}

type RecordTypeData struct {
	// Type of each member in the struct
	Name    string
	Package string
	Bases   []*Type
	Fields  []*Field
}

func (td *RecordTypeData) NumFields() int {
	return len(td.Fields)
}

func (td *RecordTypeData) NumBases() int {
	return len(td.Bases)
}

type Field struct {
	Name string
	Type *Type
}

type FunctionTypeData struct {
	// Types of the input parameters
	InputTypes []*Type

	// Types of the output parameters
	OutputTypes []*Type

	// Types of possible exceptions that can be thrown (not supported in all languages)
	ExceptionTypes []*Type
}

func (td *FunctionTypeData) NumInputs() int {
	return len(td.InputTypes)
}

func (td *FunctionTypeData) NumOutputs() int {
	return len(td.OutputTypes)
}

func (td *FunctionTypeData) NumExceptions() int {
	return len(td.ExceptionTypes)
}
