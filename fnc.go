package fanet

type FNCCommand struct {
	AircraftType int
	LiveTracking bool
	GroundType   int
}

func parseFNCCommand(tok *tokenizer) (*FNCCommand, error) {
	var c FNCCommand
	c.AircraftType = tok.hex()
	c.LiveTracking = tok.commaBool()
	c.GroundType = tok.commaHex()
	tok.endOfData()
	return &c, tok.err()
}

func (c *FNCCommand) Sentence() string {
	b := newSentenceBuilder("#FNC ")
	b.hex(c.AircraftType)
	b.commaBool(c.LiveTracking)
	b.commaHex(c.GroundType)
	b.newline()
	return b.String()
}
