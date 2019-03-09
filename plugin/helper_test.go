package plugin

import (
	"errors"
	"sync"
)

// TestingBinaryMessenger implements the BinaryMessenger interface for testing
//  purposes. It can be used as a backend in tests for BasicMessageChannel and
// StandardMethodChannel.
type TestingBinaryMessenger struct {
	channelHandlersLock sync.Mutex
	channelHandlers     map[string]BinaryMessageHandler

	// handlers mocking the other side of the BinaryMessenger
	mockChannelHandlersLock sync.Mutex
	mockChannelHandlers     map[string]BinaryMessageHandler
}

func NewTestingBinaryMessenger() *TestingBinaryMessenger {
	return &TestingBinaryMessenger{
		channelHandlers:     make(map[string]BinaryMessageHandler),
		mockChannelHandlers: make(map[string]BinaryMessageHandler),
	}
}

var _ BinaryMessenger = &TestingBinaryMessenger{} // compile-time type check

// Send sends the bytes onto the given channel.
// In this testing implementation of a BinaryMessenger, the handler for the
// channel may be set using MockSetMessageHandler
func (t *TestingBinaryMessenger) Send(channel string, message []byte) (reply []byte, err error) {
	t.mockChannelHandlersLock.Lock()
	handler := t.mockChannelHandlers[channel]
	t.mockChannelHandlersLock.Unlock()
	if handler == nil {
		return nil, errors.New("no handler set") // TODO: should this be a mock error?
	}
	return handler.HandleMessage(message)
}

// SetMessageHandler registers a binary message handler on given channel.
// In this testing implementation of a BinaryMessenger, the handler may be
// executed by calling MockSend(..).
func (t *TestingBinaryMessenger) SetMessageHandler(channel string, handler BinaryMessageHandler) {
	t.channelHandlersLock.Lock()
	t.channelHandlers[channel] = handler
	t.channelHandlersLock.Unlock()
}

// MockSend imitates a send method call from the other end of the binary
// messenger. It calls a method that was registered through SetMessageHandler.
func (t *TestingBinaryMessenger) MockSend(channel string, message []byte) (reply []byte, err error) {
	t.channelHandlersLock.Lock()
	handler := t.channelHandlers[channel]
	t.channelHandlersLock.Unlock()
	if handler == nil {
		return nil, errors.New("no handler set") // TODO: should this be a mock error?
	}
	return handler.HandleMessage(message)
}

// MockSetMessageHandler imitates a handler set at the other end of the
// binary messenger.
func (t *TestingBinaryMessenger) MockSetMessageHandler(channel string, handler BinaryMessageHandler) {
	t.mockChannelHandlersLock.Lock()
	t.mockChannelHandlers[channel] = handler
	t.mockChannelHandlersLock.Unlock()
}
