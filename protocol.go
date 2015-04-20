package bridge

type TransportRequest interface {
}

/**
 * Generates the code required to write an operation and its parameters for a
 * particular protocol.
 */
type Protocol interface {
	// Starts an operation
	EmitSendRequestMethod(opName string, opType *FunctionTypeData, argPrefix string) error

	// Starts a parameter at a given index
	EmitObjectWriterCall(argName string, argType *Type) error
	StartWritingList(childType *Type) error
	EndWritingList(childType *Type) error
	StartWritingSet(childType *Type) error
	EndWritingSet(childType *Type) error
	StartWritingMap(keyType *Type, valueType *Type) error
	EndWritingMap(keyType *Type, valueType *Type) error
	StartWritingChild(key interface{})
	EndWritingChild(key interface{})

	// Starts an operation
	EmitReadResponseMethod(opName string, opType *FunctionTypeData, argPrefix string) error

	// Starts a parameter at a given index
	EmitObjectReaderCall(argName string, argType *Type) error
	StartReadingList(childType *Type) error
	EndReadingList(childType *Type) error
	StartReadingSet(childType *Type) error
	EndReadingSet(childType *Type) error
	StartReadingMap(keyType *Type, valueType *Type) error
	EndReadingMap(keyType *Type, valueType *Type) error
	StartReadingChild(key interface{})
	EndReadingChild(key interface{})
}
