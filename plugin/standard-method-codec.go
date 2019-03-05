package plugin

import (
	"bytes"

	"github.com/pkg/errors"
)

// The first byte in a standard method envelope determines it's type.
const (
	standardMethodEnvelope_success = 0
	standardMethodEnvelope_error   = 1
)

// StandardMethodCodec implements a MethodCodec using the Flutter standard
// binary encoding.
//
// This codec tries to stay compatible with the corresponding
// StandardMethodCodec on the Dart side.
// See https://docs.flutter.io/flutter/services/StandardMethodCodec-class.html
//
// Values supported as method arguments and result payloads are those supported
// by StandardMessageCodec.
type StandardMethodCodec struct {
	// Setting a custom/extended StandardMessageCodec is not supported.
	codec StandardMessageCodec
}

var _ MethodCodec = StandardMethodCodec{}

// EncodeMethodCall fulfils the MethodCodec interface.
func (s StandardMethodCodec) EncodeMethodCall(methodCall MethodCall) (data []byte, err error) {
	buf := &bytes.Buffer{}
	err = s.codec.writeValue(buf, methodCall.method)
	if err != nil {
		return nil, errors.Wrap(err, "failed writing methodcall method name")
	}
	err = s.codec.writeValue(buf, methodCall.arguments)
	if err != nil {
		return nil, errors.Wrap(err, "failed writing methodcall arguments")
	}
	return buf.Bytes(), nil
}

// DecodeMethodCall fulfils the MethodCodec interface.
func (s StandardMethodCodec) DecodeMethodCall(data []byte) (methodCall MethodCall, err error) {
	method, err := s.codec.DecodeMessage(data)
	if err != nil {
		return methodCall, errors.Wrap(err, "failed to decode method name")
	}
	var ok bool
	methodCall.method, ok = method.(string)
	if !ok {
		return methodCall, errors.New("decoded method name is not a string")
	}
	methodCall.arguments, err = s.codec.DecodeMessage(data)
	if err != nil {
		return methodCall, errors.Wrap(err, "failed decoding method arguments")
	}
	return methodCall, nil
}

// EncodeSuccessEnvelope fulfils the MethodCodec interface.
func (s StandardMethodCodec) EncodeSuccessEnvelope(result interface{}) (data []byte, err error) {
	buf := &bytes.Buffer{}
	err = buf.WriteByte(standardMethodEnvelope_success)
	if err != nil {
		return nil, err
	}
	err = s.codec.writeValue(buf, result)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodeErrorEnvelope fulfils the MethodCodec interface.
func (s StandardMethodCodec) EncodeErrorEnvelope(code string, message string, details interface{}) (data []byte, err error) {
	buf := &bytes.Buffer{}
	err = buf.WriteByte(standardMethodEnvelope_error)
	if err != nil {
		return nil, err
	}
	err = s.codec.writeValue(buf, code)
	if err != nil {
		return nil, err
	}
	err = s.codec.writeValue(buf, message)
	if err != nil {
		return nil, err
	}
	err = s.codec.writeValue(buf, details)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecodeEnvelope fulfils the MethodCodec interface.
func (s StandardMethodCodec) DecodeEnvelope(envelope []byte) (result interface{}, err error) {
	buf := bytes.NewBuffer(envelope)
	flag, err := buf.ReadByte()
	if err != nil {
		return nil, errors.Wrap(err, "failed reading envelope flag")
	}
	switch flag {
	case standardMethodEnvelope_success:
		result, err = s.codec.readValue(buf)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode result")
		}
		return result, nil

	case standardMethodEnvelope_error:
		ferr := FlutterError{}
		var ok bool
		code, err := s.codec.readValue(buf)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode error code")
		}
		ferr.Code, ok = code.(string)
		if !ok {
			return nil, errors.New("decoded error code is not a string")
		}
		message, err := s.codec.readValue(buf)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode error message")
		}
		if message != nil {
			ferr.Message, ok = message.(string)
			if !ok {
				return nil, errors.New("decoded error message is not a string")
			}
		}
		ferr.Details, err = s.codec.readValue(buf)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode error details")
		}
		return nil, ferr
	default:
		return nil, errors.New("unknown envelope flag")
	}
}
