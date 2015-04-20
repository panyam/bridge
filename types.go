package bridge

const (
	NullType = iota
	LazyType
	BasicType
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
	Name   string
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

func (t *Type) Signature() string {
	switch t.TypeClass {
	case NullType:
		return ""
	case LazyType:
		return t.TypeData.(string)
	case BasicType:
		return t.TypeData.(string)
	case AliasType:
		return t.TypeData.(string)
	case ReferenceType:
		return "*" + t.TypeData.(*ReferenceTypeData).TargetType.Signature()
	case RecordType:
		return t.TypeData.(*RecordTypeData).Name
	case TupleType:
		out := "("
		for index, childType := range t.TypeData.(*TupleTypeData).SubTypes {
			if index > 0 {
				out += ", "
			}
			out += childType.Signature()
		}
		return out
	case FunctionType:
		funcTypeData := t.TypeData.(*FunctionTypeData)
		out := "func"
		out += TypeListSignature(funcTypeData.InputTypes)
		if funcTypeData.OutputTypes != nil {
			out += TypeListSignature(funcTypeData.OutputTypes)
		}
		if funcTypeData.ExceptionTypes != nil {
			out += " throws" + TypeListSignature(funcTypeData.ExceptionTypes)
		}
		return out
	case ListType:
		return "[]" + t.TypeData.(*ListTypeData).TargetType.Signature()
	case MapType:
		mapTypeData := t.TypeData.(*MapTypeData)
		return "map[" + mapTypeData.KeyType.Signature() + "]" + mapTypeData.ValueType.Signature()
	}
	return ""
}

func TypeListSignature(types []*Type) string {
	out := " ("
	if types != nil {
		for index, inType := range types {
			if index > 0 {
				out += ","
			}
			out += inType.Signature()
		}
	}
	out += ")"
	return out
}
