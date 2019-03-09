package plugin

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	. "github.com/stretchr/testify/assert"
)

func TestBasicMethodChannelSend(t *testing.T) {
	codec := StringCodec{}
	messenger := NewTestingBinaryMessenger()
	messenger.MockSetMessageHandler("ch", BinaryMessageHandlerFunc(func(encodedMessage []byte) ([]byte, error) {
		message, err := codec.DecodeMessage(encodedMessage)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode message")
		}
		messageString, ok := message.(string)
		if !ok {
			return nil, errors.New("message is invalid type, expected string")
		}
		reply := messageString + " world"
		encodedReply, err := codec.EncodeMessage(reply)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode message")
		}
		return encodedReply, nil
	}))
	channel := NewBasicMessageChannel(messenger, "ch", codec)
	reply, err := channel.Send("hello")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(spew.Sdump(reply))
	replyString, ok := reply.(string)
	if !ok {
		t.Fatal("reply is invalid type, expected string")
	}
	Equal(t, "hello world", replyString)
}

func TestBasicMethodChannelHandle(t *testing.T) {
	codec := StringCodec{}
	messenger := NewTestingBinaryMessenger()
	channel := NewBasicMessageChannel(messenger, "ch", codec)
	channel.HandleFunc(func(message interface{}) (reply interface{}, err error) {
		messageString, ok := message.(string)
		if !ok {
			return nil, errors.New("message is invalid type, expected string")
		}
		reply = messageString + " world"
		return reply, nil
	})
	encodedMessage, err := codec.EncodeMessage("hello")
	if err != nil {
		t.Fatalf("failed to encode message: %v", err)
	}
	encodedReply, err := messenger.MockSend("ch", encodedMessage)
	if err != nil {
		t.Fatal(err)
	}
	reply, err := codec.DecodeMessage(encodedReply)
	if err != nil {
		t.Fatalf("failed to decode reply: %v", err)
	}
	t.Log(spew.Sdump(reply))
	replyString, ok := reply.(string)
	if !ok {
		t.Fatal("reply is invalid type, expected string")
	}
	Equal(t, "hello world", replyString)
}
