package fanet

import "strings"

// A DBGCommand is a DBG command.
type DBGCommand struct {
	Zones    []string
	Severity string
}

func parseDBGCommand(tok *tokenizer) (*DBGCommand, error) {
	var c DBGCommand
	c.Zones = strings.Split(tok.string(), "|")
	c.Severity = tok.commaString()
	tok.endOfData()
	return &c, tok.err()
}

func (c *DBGCommand) Sentence() string {
	b := newSentenceBuilder("#DBG ")
	b.string(strings.Join(c.Zones, "|"))
	b.commaString(c.Severity)
	b.newline()
	return b.String()
}
