package bridge

import (
	"reflect"
)

/**
 * A InputBinder implementation extracts the request parameters out of a
 * transport object (eg http, websocket etc) and creates a ServiceOperation
 * instance that can be invoked to get the result of the service call.
 */
type InputBinder interface {
	ExtractInput(transport interface{}) (*ServiceOperation, error)
}

/**
 * OutputPresenters take the result (or error) of a service operation invocation
 * and present them back on the transport.
 */
type OutputPresenter interface {
	PresentOutput(transport interface{}, result interface{}, err error)
}

/**
 * A service operation object encapsulates the actual method (in a service) to
 * be invoked along with its arguments.
 * When the method is invoked its results are returned.  Without loss of
 * generality, requests are assumed to accept a single "request" object.
 * This ensures that the method signature itself is not complex and can compose
 * all request related parameters.
 */
type ServiceOperation struct {
	Method       reflect.Value
	RequestParam interface{}
}

func (*ServiceOperation) Invoke() (interface{}, error) {
	return nil, nil
}
