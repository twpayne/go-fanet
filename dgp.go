package fanet

type DGPCommand struct {
	PowerMode bool
}

func parseDGPCommand(tok *tokenizer) (*DGPCommand, error) {
	var c DGPCommand
	c.PowerMode = tok.bool()
	tok.endOfData()
	return &c, tok.err()
}

func (c *DGPCommand) Sentence() string {
	b := newSentenceBuilder("#DGP ")
	b.bool(c.PowerMode)
	b.newline()
	return b.String()
}
