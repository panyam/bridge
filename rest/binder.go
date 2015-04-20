package rest

import (
	// "github.com/gorilla/mux"
	"github.com/panyam/relay/bindings"
	"github.com/panyam/relay/utils"
	"log"
	"net/http"
	"reflect"
)

type HttpBinding struct {
	/**
	 * Methods required to match for this binding to get triggered.
	 */
	Methods []string

	/**
	 * The URL which will trigger this binding.
	 */
	Url string

	// Mappings between a query or BODY parameter to a key with the request
	ParamMappings map[string][]string

	// Mappings between a path variable to a key with the request
	VarMappings map[string][]string

	/**
	 * The service that needs to be invoked when the binding matches.
	 */
	Service interface{}

	/**
	 * Name of the operation to invoke.
	 */
	Operation string

	/**
	 * Type of the request object to be created and populated.
	 */
	RequestType      reflect.Type
	RequestTypeIsPtr bool

	/**
	 * The method corresponding to the operation within the service.
	 */
	Method reflect.Value
}

func NewHttpBinding(url string, methods []string, service interface{}, operation string) *HttpBinding {
	out := HttpBinding{Url: url, Methods: methods, Service: service, Operation: operation}
	out.Method = utils.GetMethod(service, operation)
	out.RequestTypeIsPtr, out.RequestType = utils.GetParamType(out.Method, 0)
	if out.RequestType != nil {
		return &out
	}
	log.Println("Invalid operation on service: ", operation)
	return nil
}

/**
 * Bindings are based on:
 * Body Parser:
 * 	JSON, XML - Request type maps directly to content
 * 	URL Query Params and FormData - Maps to specific parameters in the request
 *
 * 	In general the algo is parse the body first (if any) and create the request
 * 	object.  Then for the values not found in the body get them from the query
 * 	parameters and URL parameters.
 *
 * 	By default, The names of fields in JSON, XML or Form data will map directly
 * 	to the request attributes.   With Form data container data wont be possible
 * 	but with JSON/XML they will be extracted if available.  Query parameters
 * 	will also by default map directly to request field names.  However these can
 * 	be overridden by a mapping of the form:
 *
 * 	NameMappings = {
 * 		"FormParamName1" = "Request.Field1"
 * 		"FormParamName2" = "Request.Field1.ChildField"
 * 	}
 *
 * 	With parameterized URL paths, URLs will be specified as:
 *
 * 	/path1/path2/{param1:Request.Field1}/path3/{param2:Request.Field2}/
 */
func (hb *HttpBinding) ExtractRequest(request *http.Request) (*bindings.ServiceOperation, error) {
	param := reflect.New(hb.RequestType)
	out := bindings.ServiceOperation{Method: hb.Method, RequestParam: param}

	// variables = mux.Vars(request)
	return &out, nil
}

type HttpInputBinder struct {
	// Methods that are ok for this
	Bindings []*HttpBinding
}

func (h *HttpInputBinder) ExtractInput(transportRequest interface{}) (*bindings.ServiceOperation, error) {
	request := transportRequest.(*http.Request)
	binding := h.MatchBinding(request)
	// then extract the operation items
	if binding == nil {
		// No binding found so return
		return nil, nil
	}

	// extract request from binding
	return binding.ExtractRequest(request)
}

func (h *HttpInputBinder) MatchBinding(request *http.Request) *HttpBinding {
	// TODO: Use a Trie to store matches based on prefixes
	for _, binding := range h.Bindings {
		matched := binding.Methods == nil
		if !matched {
			for _, method := range binding.Methods {
				if method == request.Method {
					matched = true
					break
				}
			}
		}

		if !matched {
			return nil
		}

		// now see if URL matches
		// variables = mux.Vars(request)
	}
	return nil
}

func (h *HttpInputBinder) AddBinding(binding *HttpBinding) {
	// TODO: store the bindings as Trie so we dont have to search every
	// match each time
	h.Bindings = append(h.Bindings, binding)
}
