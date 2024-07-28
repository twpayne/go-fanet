package fanet

type FNMCommand struct {
	Mode int
}

func parseFNMCommand(tok *tokenizer) (*FNMCommand, error) {
	var c FNMCommand
	c.Mode = tok.int()
	tok.endOfData()
	return &c, tok.err()
}

func (c *FNMCommand) Sentence() string {
	b := newSentenceBuilder("#FNM ")
	b.int(c.Mode)
	b.newline()
	return b.String()
}
