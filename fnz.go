package fanet

// An FNZResponse is an FNZ response.
type FNZResponse struct {
	Zone     int
	ZoneName string
}

func parseFNZResponse(tok *tokenizer) (*FNZResponse, error) {
	var fnz FNZResponse
	fnz.Zone = tok.int()
	fnz.ZoneName = tok.commaString()
	tok.endOfData()
	return &fnz, tok.err()
}

func (r *FNZResponse) Address() string {
	return "FNZ"
}
