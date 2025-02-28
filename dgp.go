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

type DGPResponse struct {
	PowerMode bool
}

func parseDGPResponse(tok *tokenizer) (*DGPResponse, error) {
	var r DGPResponse
	r.PowerMode = tok.int() != 0
	tok.endOfData()
	return &r, tok.err()
}

func (r *DGPResponse) Address() string {
	return "DGP"
}
