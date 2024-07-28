package fanet

// A DGJCommand is a DGJ command.
type DGJCommand struct {
	Bootloader string
}

func parseDGJCommand(tok *tokenizer) (*DGJCommand, error) {
	var c DGJCommand
	c.Bootloader = tok.string()
	tok.endOfData()
	return &c, tok.err()
}

func (c *DGJCommand) Sentence() string {
	b := newSentenceBuilder("#DGJ ")
	b.string(c.Bootloader)
	b.newline()
	return b.String()
}
