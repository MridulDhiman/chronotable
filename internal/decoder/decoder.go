package decoder

import (
	"encoding/gob"
	"io"
)

type Decoder struct {
	reader *gob.Decoder
}

func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{
		reader: gob.NewDecoder(reader),
	}
}

func (decoder *Decoder) Decode(e any) error {
	return decoder.reader.Decode(e)
}