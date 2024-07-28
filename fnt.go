package fanet

import (
	"fmt"
)

// An FNTCommand is an FNT request.
type FNTCommand struct {
	Type        int
	Destination ID
	Forward     bool
	AckRequired bool
	Payload     []byte
	Signature   Optional[int]
}

func parseFNTCommand(tok *tokenizer) (*FNTCommand, error) {
	var c FNTCommand
	c.Type = tok.hex()
	c.Destination.Manufacturer = tok.commaHex()
	c.Destination.Device = tok.commaHex()
	c.Forward = tok.commaBool()
	c.AckRequired = tok.commaBool()
	payloadLength := tok.commaHex()
	c.Payload = tok.commaHexBytes()
	if len(c.Payload) != payloadLength {
		return nil, tok.errOr(fmt.Errorf("payload length mismatch, want %d, got %d", payloadLength, len(c.Payload)))
	}
	if !tok.atEndOfData() {
		c.Signature = NewOptional(tok.commaHex())
	}
	tok.endOfData()
	return &c, tok.err()
}

func (c *FNTCommand) Sentence() string {
	b := newSentenceBuilder("#FNT ")
	b.hex(c.Type)
	b.commaHex(c.Destination.Manufacturer)
	b.commaHex(c.Destination.Device)
	b.commaBool(c.Forward)
	b.commaBool(c.AckRequired)
	b.commaHex(len(c.Payload))
	b.commaBytes(c.Payload)
	if c.Signature.Valid {
		b.commaHex(c.Signature.Value)
	}
	b.newline()
	return b.String()
}
