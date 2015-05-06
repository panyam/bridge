package bridge

import (
	"fmt"
)

const (
	NullType = iota
	UnresolvedType
	NamedType
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

func (t *Type) TypeClassString() string {
	switch t.TypeClass {
	case NullType:
		return "NullType"
	case UnresolvedType:
		return "UnresolvedType"
	case NamedType:
		return "NamedType"
	case AliasType:
		return "AliasType"
	case ReferenceType:
		return "ReferenceType"
	case TupleType:
		return "TupleType"
	case RecordType:
		return "RecordType"
	case FunctionType:
		return "FunctionType"
	case ListType:
		return "ListType"
	case MapType:
		return "MapType"
	}
	return ""
}

/**
 * Tells if the value will be passed by value or by reference.
 */
func (t *Type) IsValueType() bool {
	if t.TypeClass == ReferenceType || t.TypeClass == ListType ||
		t.TypeClass == MapType || t.TypeClass == FunctionType {
		return false
	} else if t.TypeClass == AliasType {
		return t.AsAliasType().TargetType.IsValueType()
	} else if t.TypeClass == NamedType || t.TypeClass == RecordType {
		return true
	}
	return false
}

func (t *Type) IsNullType() bool       { return t.TypeClass == NullType }
func (t *Type) IsUnresolvedType() bool { return t.TypeClass == UnresolvedType }
func (t *Type) IsNamedType() bool      { return t.TypeClass == NamedType }
func (t *Type) IsAliasType() bool      { return t.TypeClass == AliasType }
func (t *Type) IsReferenceType() bool  { return t.TypeClass == ReferenceType }
func (t *Type) IsTupleType() bool      { return t.TypeClass == TupleType }
func (t *Type) IsRecordType() bool     { return t.TypeClass == RecordType }
func (t *Type) IsFunctionType() bool   { return t.TypeClass == FunctionType }
func (t *Type) IsListType() bool       { return t.TypeClass == ListType }
func (t *Type) IsMapType() bool        { return t.TypeClass == MapType }

func (t *Type) AsNamedType() *NamedTypeData         { return t.TypeData.(*NamedTypeData) }
func (t *Type) AsAliasType() *AliasTypeData         { return t.TypeData.(*AliasTypeData) }
func (t *Type) AsReferenceType() *ReferenceTypeData { return t.TypeData.(*ReferenceTypeData) }
func (t *Type) AsTupleType() *TupleTypeData         { return t.TypeData.(*TupleTypeData) }
func (t *Type) AsRecordType() *RecordTypeData       { return t.TypeData.(*RecordTypeData) }
func (t *Type) AsFunctionType() *FunctionTypeData   { return t.TypeData.(*FunctionTypeData) }
func (t *Type) AsListType() *ListTypeData           { return t.TypeData.(*ListTypeData) }
func (t *Type) AsMapType() *MapTypeData             { return t.TypeData.(*MapTypeData) }

func (t *Type) LeafType() *NamedTypeData {
	switch typeData := t.TypeData.(type) {
	case *NamedTypeData:
		return typeData
	case *AliasTypeData:
		return typeData.TargetType.LeafType()
	case *ReferenceTypeData:
		return typeData.TargetType.LeafType()
	case *RecordTypeData:
		return &typeData.NamedTypeData
	}
	return nil
}

func (t *Type) String() string {
	return fmt.Sprintf("{%d - %s}", t.TypeClass, t.TypeData)
}

func NewType(typeCls int, data interface{}) *Type {
	return &Type{TypeClass: typeCls, TypeData: data}
}

type NamedTypeData struct {
	Name    string
	Package string
}

type AliasTypeData struct {
	// Type this is an alias/typedef for
	Name       string
	TargetType *Type
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
	NamedTypeData

	// Type of each member in the struct
	Bases  []*Type
	Fields []*Field
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
