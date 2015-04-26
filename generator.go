package bridge

import (
	"io"
)

/**
 * Responsible for generating the code for the client classes.
 */
type Generator interface {
	/**
	 * Emits the class that acts as a client for the service.
	 */
	EmitClientClass(writer io.Writer, serviceType *Type) error

	/**
	 * For a given service operation, emits a method which:
	 * 1. Has inputs the same as those of the underlying service operation,
	 * 2. creates a transport level request
	 * 3. Sends the transport level request
	 * 4. Gets a response from the transport level and returns it
	 */
	EmitServiceCallMethod(writer io.Writer, opName string, opType *FunctionTypeData, argPrefix string) error

	/**
	 * Emits the writer for a particular type.
	 */
	EmitTypeWriter(writer io.Writer, argType *Type) error

	/**
	 * Emits the reader for a particular type.
	 */
	EmitTypeReader(writer io.Writer, argType *Type) error
}
