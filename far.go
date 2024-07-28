package fanet

import "fmt"

// A FARResponse is a FAR response.
type FARResponse struct {
	Type      string
	Status    int
	StatusStr string
}

func parseFARResponse(tok *tokenizer) (*FARResponse, error) {
	var r FARResponse
	r.Type = tok.string()
	switch r.Type {
	case "OK":
	case "ERR":
		r.Status = tok.commaInt()
		r.StatusStr = tok.commaString()
	default:
		return nil, fmt.Errorf("%s: unknown FAR type", r.Type)
	}
	tok.endOfData()
	return &r, tok.err()
}

func (r *FARResponse) Err() error {
	if r.Type == "OK" {
		return nil
	}
	return fmt.Errorf("%d: %s", r.Status, r.StatusStr)
}
