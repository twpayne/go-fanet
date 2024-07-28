package fanet

import (
	"time"
)

// A FATCommand is a FAT command.
type FATCommand struct {
	Delay time.Duration
}

func parseFATCommand(tok *tokenizer) (*FATCommand, error) {
	var c FATCommand
	c.Delay = time.Duration(tok.int()) * time.Millisecond
	tok.endOfData()
	return &c, tok.err()
}

func (c *FATCommand) Sentence() string {
	b := newSentenceBuilder("#FAT ")
	b.int(int(c.Delay / time.Millisecond))
	b.newline()
	return b.String()
}
