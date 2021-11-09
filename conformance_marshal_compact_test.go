package main

import (
	"bytes"
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
	. "github.com/onsi/gomega"
)

func TestMarshalling(t *testing.T) {
	g := NewWithT(t)

	pd := partialDoc{
		fields: map[string][]byte{
			"a": []byte(`"aValue"`),
			"b": []byte(`"bValue"`),
		},
	}

	out, _ := json.Marshal(&pd)
	// encoding/json.marshalerEncoder(e *encodeState, v reflect.Value, opts encOpts) is running
	// compact after calling the custom MarshalJSON func

	jsonIterOut, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&pd)


	g.Expect(out).To(Equal(jsonIterOut))
}

type partialDoc struct {
	fields map[string][]byte
}

func (n *partialDoc) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if _, err := buf.WriteString("{"); err != nil {
		return nil, err
	}

	i := 0
	for k, v := range n.fields {
		if i > 0 {
			if _, err := buf.WriteString(", "); err != nil {
				return nil, err
			}
		}
		key, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		if _, err := buf.Write(key); err != nil {
			return nil, err
		}

		if _, err := buf.WriteString(": "); err != nil {
			return nil, err
		}

		if _, err := buf.Write(v); err != nil {
			return nil, err
		}
		i++
	}
	if _, err := buf.WriteString("}"); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
