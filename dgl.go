package fanet

// A DGLCommand is a DGL command.
type DGLCommand struct {
	Frequency int
	DBm       int
}

func parseDGLCommand(tok *tokenizer) (*DGLCommand, error) {
	var c DGLCommand
	c.Frequency = tok.int()
	c.DBm = tok.commaInt()
	tok.endOfData()
	return &c, tok.err()
}

func (c *DGLCommand) Sentence() string {
	b := newSentenceBuilder("#DGL ")
	b.int(c.Frequency)
	b.commaInt(c.DBm)
	b.newline()
	return b.String()
}
