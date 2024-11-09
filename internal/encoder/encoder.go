package encoder

import (
	"encoding/gob"
	"io"
)

type Encoder struct {
	writer *gob.Encoder
}

func NewEncoder(writer io.Writer) *Encoder {
	return &Encoder{
		writer: gob.NewEncoder(writer),
	}
}

func (encoder *Encoder) Encode(e any) error {
	return encoder.writer.Encode(e)
}
