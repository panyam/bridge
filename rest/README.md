
The binding object should be used by both the server as well as the client
generator.   This way we dont have repetitions on building this.  So one
binding is a Restful HTTP binding (others are RCP etc).
A binding should specify the service, operation and details about the
transport level endpoint.  eg:


```
Binding:
	Service: IUserService,
	Operation: CreateTeam,
	Input: 
		Params: request
		Endpoint: POST /teams/
	Output:
		Params: request, error
		Body: // Only one of the following is required
			Template: A template string to which the outputs are passed as an array to be rendered.
			TemplateFile: Name of the template file to which the outputs are passed as an array to be rendered.
			Presenter: A presenter function that presents the output to the
			transport
```

The above binding is for creating a team.  It specifies the endpoint on both
sides.  The type of the input will be inferred by inspecting the parameters
for the service method.  Variable bindings are specified with {...}.  Input
parameters can be given names as above.  The input endpoint and the body will
all the details to read the request object.  The input endpoint path can have
parametrised values to denote which values in the request are to be set.

Similary output presenter has its own section.  

This binding should have enough information to allow the following:

	1. Generation of request parsers
	2. Generation of request presenters
	3. Generation of response parsers
	4. Generation of response presenters

Note that this is specific to one kind of transport (http).  The input
and output structures for a service could be different for other kinds of
transport (eg binary or socket or RPC etc)

How will all this be integrated back?

1. go generate httpbindings -input bindings.go -lang <language> -o outfile -- <language specific options>

This will generate a list of files (or a single file depending on language) that
has serialiser and deserializer of service request objects to and from http
(taking into account content types and accept headers)

2. This will also generate http handler functions for each service operation, eg:

```
type ITeamServiceHandler struct {
	teamService ITeamService
	RequestDecorator func(req *http.Request) (*http.Request, error)
}

func (svc ITeamServiceClient) SendCreateTeam(req *CreateTeamRequest) (*CreateTeamResponse, error) {
	httpreq, err := SerializeCreateTeamRequest(req)
	if err != nil {
		return nil, err
	}
	if svc.RequestDecorator != nil {
		httpreq, err = svc.RequestDecorator(httpreq)
		if err != nil {
			return nil, err
		}
	}
	resp, err := SendHttpRequest(httpreq)
	if err != nil {
		return nil, err
	}
	return DeserializeCreateTeamResponse(resp, httpreq)
}
```

What will be generated for each language is:

* Serialize<operation>Request
* Deserialize<operation>Request
* Serialize<operation>Response
* Deserialize<operation>Response
* <Service>Handler class
* <Service>Client class

