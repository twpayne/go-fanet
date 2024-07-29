package fanet

import "fmt"

// An FNRResponse is an FNR response.
type FNRResponse struct {
	Type        string
	Status      int
	StatusStr   string
	Destination ID
}

func parseFNRResponse(tok *tokenizer) (*FNRResponse, error) {
	var r FNRResponse
	r.Type = tok.string()
	switch r.Type {
	case "OK":
	case "ERR", "WRN", "MSG":
		r.Status = tok.commaInt()
		r.StatusStr = tok.commaString()
	case "ACK", "NACK":
		r.Destination.Manufacturer = tok.commaHex()
		r.Destination.Device = tok.commaHex()
	default:
		return nil, fmt.Errorf("%s: unknown FNR type", r.Type)
	}
	tok.endOfData()
	return &r, tok.err()
}

func (r *FNRResponse) Err() error {
	if r.Type != "ERR" {
		return nil
	}
	return fmt.Errorf("%d: %s", r.Status, r.StatusStr)
}
