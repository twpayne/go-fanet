package fanet

// An FNACommand is an FNA command.
type FNACommand struct{}

func parseFNACommand(tok *tokenizer) (*FNACommand, error) {
	var r FNACommand
	tok.endOfData()
	return &r, tok.err()
}

func (FNACommand) Sentence() string {
	return "#FNA\n"
}

// An FNAResponse is an FNA response.
type FNAResponse struct {
	ID ID
}

func parseFNAResponse(tok *tokenizer) (*FNAResponse, error) {
	var r FNAResponse
	r.ID.Manufacturer = tok.hex()
	r.ID.Device = tok.commaHex()
	tok.endOfData()
	return &r, tok.err()
}

func (r *FNAResponse) Address() string {
	return "FNA"
}
