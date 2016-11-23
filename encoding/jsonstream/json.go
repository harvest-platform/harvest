package jsonstream

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
)

var ErrClosed = errors.New("json: stream closed")

func isArray(br *bufio.Reader) (bool, error) {
	for {
		b, err := br.Peek(1)
		if err != nil {
			return false, err
		}

		switch b[0] {
		// Ignore whitespace.
		case ' ', '\n', '\t':
		case '[':
			return true, nil

		case '{':
			return false, nil

		default:
			return false, errors.New("JSON must be an object or array of objects")
		}
	}
}

type Encoder struct {
	w      io.Writer
	json   *json.Encoder
	first  bool
	closed bool
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:     w,
		json:  json.NewEncoder(w),
		first: true,
	}
}

func (e *Encoder) Encode(v interface{}) error {
	if e.closed {
		return ErrClosed
	}

	if e.first {
		e.w.Write([]byte("["))
		e.first = false
	} else {
		e.w.Write([]byte(","))
	}

	return e.json.Encode(v)
}

// Close closes the stream. This _does not_ close the underlying io.Writer.
func (e *Encoder) Close() error {
	if e.closed {
		return ErrClosed
	}

	// In case no elements are written, open the brace.
	if e.first {
		e.w.Write([]byte("["))
	}

	e.w.Write([]byte("]"))
	e.closed = true

	return nil
}

type Decoder struct {
	br    *bufio.Reader
	json  *json.Decoder
	first bool
	err   error
}

func NewDecoder(r io.Reader) *Decoder {
	br := bufio.NewReader(r)
	return &Decoder{
		br:    br,
		json:  json.NewDecoder(br),
		first: true,
	}
}

func (e *Decoder) init() {
	if e.first {
		e.first = false

		ok, err := isArray(e.br)
		if err != nil {
			e.err = err
		}

		// If this is an array, read the bracket.
		if ok {
			if _, err := e.json.Token(); err != nil {
				e.err = err
			}
		}
	}
}

func (e *Decoder) More() bool {
	if e.first {
		e.init()
	}

	return e.json.More()
}

func (e *Decoder) Decode(v interface{}) error {
	if e.first {
		e.init()
	}

	if e.err != nil {
		return e.err
	}

	err := e.json.Decode(v)
	if err != nil {
		return err
	}

	// Read closing bracket.
	if !e.json.More() {
		if _, err := e.json.Token(); err != nil {
			return err
		}
	}

	return nil
}
