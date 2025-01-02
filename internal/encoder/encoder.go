package encoder

import (
	"encoding/gob"
	"io"
)


type IEncoder interface {
	Encode(e any) error
}

type GobEncoder struct {
	writer *gob.Encoder
}

func NewEncoder(writer io.Writer) *GobEncoder {
	return &GobEncoder{
		writer: gob.NewEncoder(writer),
	}
}

func (encoder *GobEncoder) Encode(e any) error {
	return encoder.writer.Encode(e)
}
