package fanet

type FAPCommand struct {
	PowerMode bool
}

func parseFAPCommand(tok *tokenizer) (*FAPCommand, error) {
	var c FAPCommand
	c.PowerMode = tok.bool()
	tok.endOfData()
	return &c, tok.err()
}

func (c *FAPCommand) Sentence() string {
	b := newSentenceBuilder("#FAP ")
	b.bool(c.PowerMode)
	b.newline()
	return b.String()
}
