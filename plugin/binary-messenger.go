package plugin

type BinaryMessenger interface {
	Send(channel string, message []byte) (reply []byte, err error)
	SetMessageHandler(channel string, handler BinaryMessageHandler)
}

type BinaryMessageHandler interface {
	handleMessage(message []byte) (reply []byte, err error)
}
