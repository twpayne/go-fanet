package fanet

import "fmt"

// A DGRResponse is a DGR response.
type DGRResponse struct {
	Type      string
	Status    int
	StatusStr string
}

func parseDGRResponse(tok *tokenizer) (*DGRResponse, error) {
	var r DGRResponse
	r.Type = tok.string()
	switch r.Type {
	case "OK":
	case "ERR":
		r.Status = tok.commaInt()
		r.StatusStr = tok.commaString()
	default:
		return nil, fmt.Errorf("%s: unknown DGR type", r.Type)
	}
	tok.endOfData()
	return &r, tok.err()
}

func (r *DGRResponse) Err() error {
	if r.Type == "OK" {
		return nil
	}
	return fmt.Errorf("%d: %s", r.Status, r.StatusStr)
}
