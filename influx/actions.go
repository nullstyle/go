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
