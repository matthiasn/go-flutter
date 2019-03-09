package plugin

type BinaryMessenger interface {
	Send(channel string, encodedMessage []byte) (encodedReply []byte, err error)
	SetMessageHandler(channel string, handler BinaryMessageHandler)
}

type BinaryMessageHandler interface {
	// TODO: the re-use of the term "Handle" is a bit confusing. There's just
	// too much "handle" going on... net/http used ServeHTTP as method. We need
	// something like that as well to set the method apart from the handler type.
	HandleMessage(encodedMessage []byte) (encodedReply []byte, err error)
}

// The BinaryMessageHandlerFunc type is an adapter to allow the use of
// ordinary functions as binary message handlers. If f is a function
// with the appropriate signature, BinaryMessageHandlerFunc(f) is a
// BinaryMessageHandler that calls f.
type BinaryMessageHandlerFunc func(encodedMessage []byte) (encodedReply []byte, err error)

// HandleMessage calls f(encodedMessage)
func (b BinaryMessageHandlerFunc) HandleMessage(encodedMessage []byte) (encodedReply []byte, err error) {
	return b(encodedMessage)
}
