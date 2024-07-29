package fanet

import "fmt"

// An DBRResponse is an DBR response.
type DBRResponse struct {
	Type      string
	Status    int
	StatusStr string
}

func parseDBRResponse(tok *tokenizer) (*DBRResponse, error) {
	var r DBRResponse
	r.Type = tok.string()
	switch r.Type {
	case "OK":
	case "ERR":
		r.Status = tok.commaInt()
		r.StatusStr = tok.commaString()
	default:
		return nil, fmt.Errorf("%s: unknown DBR type", r.Type)
	}
	tok.endOfData()
	return &r, tok.err()
}

func (r *DBRResponse) Err() error {
	if r.Type == "OK" {
		return nil
	}
	return fmt.Errorf("%d: %s", r.Status, r.StatusStr)
}
