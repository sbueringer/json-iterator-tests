package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
	. "github.com/onsi/gomega"
)

func TestUnmarshalErrorPropagation(t *testing.T) {
	g := NewWithT(t)

	doc := testDoc{}

	err := json.Unmarshal([]byte(`{"key": "value"}`), &doc)
	// Returns:
	// &syntaxError{}
	// err returned by testdoc.MarshalJSON is returned 1:1 in:
	// encoding/json.decodeState.object(v reflect.Value) error

	// errors.As can find the syntaxError
	se := &syntaxError{}
	ok := errors.As(err, &se)
	g.Expect(ok).To(BeTrue())

	// err is a *syntaxError
	_, ok = err.(*syntaxError)
	g.Expect(ok).To(BeTrue())

	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(`{"key": "value"}`), &doc)
	// Returns:
	// errors.errorString{
	//   s: unmarshalerDecoder: syntax error: syntax error, error found in #10 byte of ...|: "value"}|..., bigger context ...|{"key": "value"}|...
	// }
	// Error is recreated based on syntaxError.Err() in:
	// json-iterator/go.Iterator.ReportError(operation string, msg string)
	// json-iterator/go.unmarshalerDecoder.Decode(ptr unsafe.Pointer, iter *Iterator)

	// errors.As should be able to find the syntaxError
	se = &syntaxError{}
	ok = errors.As(err, &se)
	g.Expect(ok).To(BeTrue())
}

func TestMarshalErrorPropagation(t *testing.T) {
	g := NewWithT(t)

	doc := testDoc{}

	_, err := json.Marshal(&doc)
	// Returns:
	// &json.MarshalerError{
	//   Err: &syntaxError{},
	// }
	// err returned by testdoc.MarshalJSON is wrapped in a MarshalerError in:
	// encoding/json.marshalerEncoder(e *encodeState, v reflect.Value, opts encOpts)

	// errors.As can find the syntaxError
	se := &syntaxError{}
	ok := errors.As(err, &se)
	g.Expect(ok).To(BeTrue())

	// err is not a *syntaxError
	_, ok = err.(*syntaxError)
	g.Expect(ok).To(BeFalse())

	_, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&doc)
	// Returns:
	// &syntaxError{}
	// err returned by testdoc.MarshalJSON is returned 1:1 in:
	// json-iterator/go.marshalerEncoder.Encode(ptr unsafe.Pointer, stream *Stream)

	// errors.As can find the syntaxError
	se = &syntaxError{}
	ok = errors.As(err, &se)
	g.Expect(ok).To(BeTrue())

	// err is a *syntaxError
	_, ok = err.(*syntaxError)
	g.Expect(ok).To(BeTrue())
}

type testDoc struct{}

type syntaxError struct {
	msg string
}

func (n *testDoc) MarshalJSON() ([]byte, error) {
	return nil, &syntaxError{
		msg: "syntax error",
	}
}

func (n *testDoc) UnmarshalJSON(_ []byte) error {
	return &syntaxError{
		msg: "syntax error",
	}
}

func (err *syntaxError) Error() string {
	return fmt.Sprintf("syntax error: %s", err.msg)
}
