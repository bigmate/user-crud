package user

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

type modelEncoder struct {
	buf *bytes.Buffer
	err error
}

// NewModelEncoder encodes data as json and returns
// io.Reader of encoded data.
// It's very convenient to modify encoding here if you want
// to send data in other formats via notifier service
func NewModelEncoder(data interface{}) io.Reader {
	mdEnc := &modelEncoder{
		buf: &bytes.Buffer{},
	}
	enc := json.NewEncoder(mdEnc.buf)

	if err := enc.Encode(data); err != nil {
		log.Printf("modelencoder: failed to encode data: %v\n", err)
		mdEnc.err = err
	}

	return mdEnc
}

func (m *modelEncoder) Read(p []byte) (n int, err error) {
	if m.err != nil {
		return 0, err
	}

	return m.buf.Read(p)
}
