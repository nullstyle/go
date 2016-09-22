package influx

var (
	// StateLoaded is a lifecycle action that is fired whenever the influx.Store
	// has initialized or loaded the state tracked by the store.  It occurs once
	// at store creation time, and then subsequently after any successful load
	// operation.  Handlers should respond to this action by initializing the
	// unexported state based upon any input fields or context.
	StateLoaded Action = stateLoaded

	// StateWillSave is a lifecycle action that is fired when the influx.Store
	// will serialize the state to a snapshot. Handler implementations should use
	// this opportunity to clean themselves up in any way necessary for a
	// successful saved state to be created.
	StateWillSave Action = stateWillSave
)

// Request is one of the fundamental methods of inter-component communication.
// It is used to build request/response style patterns of communication within
// influx.  As an example, see HttpClient and HttpResponse within this package.
// Components can initiate an http request using the client, saving the returned
// request id in their embedded HttpResponse component.
type Request struct {
	// ID a unique identifier for request. Unique is application defined.  TODO:
	// change to a type backed by context.Context.
	ID int

	// Unprotected is true if if the request should be considered to contain
	// low-sensitivity data.  Any request components that store a request can use
	// the value to inform their serialization strategy.
	Unprotected bool
}
