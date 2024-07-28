package fanet

// A DGVCommand is a DGV command.
type DGVCommand struct{}

func parseDGVCommand(tok *tokenizer) (*DGVCommand, error) {
	var r DGVCommand
	tok.endOfData()
	return &r, tok.err()
}

func (DGVCommand) Sentence() string {
	return "#DGV\n"
}

// A DGVResponse is a DGV response.
type DGVResponse struct {
	BuildDateCode string
}

func parseDGVResponse(tok *tokenizer) (*DGVResponse, error) {
	var r DGVResponse
	r.BuildDateCode = tok.rest()
	tok.endOfData()
	return &r, tok.err()
}
