package json

import (
	"encoding/json"
	"errors"
	"io"
)

var ErrClosed = errors.New("json: stream closed")

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
	r     io.Reader
	json  *json.Decoder
	first bool
	more  bool
	err   error
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:     r,
		json:  json.NewDecoder(r),
		first: true,
	}
}

func (e *Decoder) More() bool {
	return e.json.More()
}

func (e *Decoder) Err() error {
	return e.err
}

func (e *Decoder) Decode(v interface{}) error {
	if e.first {
		// Read opening bracket.
		_, err := e.json.Token()
		if err != nil {
			return err
		}
		e.first = false
	}

	err := e.json.Decode(v)
	if err != nil {
		return err
	}

	// Read closing bracket.
	if !e.json.More() {
		e.json.Token()
	}

	return nil
}
