package plugin

import "github.com/pkg/errors"

// BasicMessageHandler defines the interfece for a basic message handler.
type BasicMessageHandler interface {
	handleMessage(message interface{}) (reply interface{}, err error)
}

// BasicMessageChannel presents named channel for communicating with the Flutter
// application using basic, asynchronous message passing.
//
// Messages are encoded into binary before being sent, and binary messages
/// received are decoded into. The MessageCodec used must be compatible with the
// one used by the Flutter application. This can be achieved by creating a
// BasicMessageChannel counterpart of this channel on the Dart side.
// See: https://docs.flutter.io/flutter/services/BasicMessageChannel-class.html
//
// The static Go type of messages sent and received is interface{}, but only
// values supported by the specified MessageCodec can be used.
//
// The logical identity of the channel is given by its name. Identically named
// channels will interfere with each other's communication.
type BasicMessageChannel struct {
	messenger BinaryMessenger
	name      string
	codec     MessageCodec
}

// NewBasicMessageChannel creates a BasicMessageChannel.
//
// Call SetMessageHandler on the returned BasicMessageChannel to provide the
// channel with a handler for incomming messages.
func NewBasicMessageChannel(messenger BinaryMessenger, name string, codec MessageCodec) *BasicMessageChannel {
	return &BasicMessageChannel{
		messenger: messenger,
		name:      name,
		codec:     codec,
	}
}

// Send encodes and sends the specified message to the Flutter application and
// returns the reply, or an error.
func (b *BasicMessageChannel) Send(message interface{}) (reply interface{}, err error) {
	encodedMessage, err := b.codec.EncodeMessage(message)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode outgoing message")
	}
	encodedReply, err := b.messenger.Send(b.name, encodedMessage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send outgoing message")
	}
	reply, err = b.codec.DecodeMessage(encodedReply)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode incomming reply")
	}
	return reply, nil
}

// SetMessageHandler registers a message handler on this channel for receiving
// messages sent from the Flutter application.
//
// Consecutive calls override any existing handler registration for (the name
// of) this channel.
//
// When given nil as handler, any incoming message on this channel will be
// handled silently by sending a nil reply (null on the dart side).
func (b *BasicMessageChannel) SetMessageHandler(handler BasicMessageHandler) {
	b.messenger.SetMessageHandler(b.name, incommingBasicMessageHandler{
		codec:   b.codec,
		handler: handler,
	})
}

// incommingBasicMessageHandler handles binary messages using
type incommingBasicMessageHandler struct {
	codec   MessageCodec
	handler BasicMessageHandler
}

var _ BinaryMessageHandler = incommingBasicMessageHandler{} // compile-time type check

// handleMessage decodes incoming binary, calls the handler, and encodes the
// outgoing reply.
func (i incommingBasicMessageHandler) handleMessage(encodedMessage []byte) (encodedReply []byte, err error) {
	if i.handler == nil {
		return nil, nil
	}
	message, err := i.codec.DecodeMessage(encodedMessage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode incomming message")
	}
	reply, err := i.handler.handleMessage(message)
	if err != nil {
		return nil, errors.Wrap(err, "handler for incoming basic message failed")
	}
	encodedReply, err = i.codec.EncodeMessage(reply)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode outgoing reply")
	}
	return encodedReply, nil
}
