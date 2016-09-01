package json

import (
	"encoding/json"
	"errors"
	"io"
)

var ErrStreamClosed = errors.New("json: stream closed")

type StreamEncoder struct {
	w      io.Writer
	json   *json.Encoder
	first  bool
	closed bool
}

func NewStreamEncoder(w io.Writer) *StreamEncoder {
	return &StreamEncoder{
		w:     w,
		json:  json.NewEncoder(w),
		first: true,
	}
}

func (e *StreamEncoder) Encode(v interface{}) error {
	if e.closed {
		return ErrStreamClosed
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
func (e *StreamEncoder) Close() error {
	if e.closed {
		return ErrStreamClosed
	}

	// In case no elements are written, open the brace.
	if e.first {
		e.w.Write([]byte("["))
	}

	e.w.Write([]byte("]"))
	e.closed = true

	return nil
}

type StreamDecoder struct {
	r     io.Reader
	json  *json.Decoder
	first bool
	more  bool
	err   error
}

func NewStreamDecoder(r io.Reader) *StreamDecoder {
	return &StreamDecoder{
		r:     r,
		json:  json.NewDecoder(r),
		first: true,
	}
}

func (e *StreamDecoder) More() bool {
	return e.json.More()
}

func (e *StreamDecoder) Err() error {
	return e.err
}

func (e *StreamDecoder) Decode(v interface{}) error {
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
