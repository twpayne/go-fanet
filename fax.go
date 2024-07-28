package fanet

import "time"

// A FAXCommand is a FAX command.
type FAXCommand struct{}

func parseFAXCommand(tok *tokenizer) (*FAXCommand, error) {
	var r FAXCommand
	tok.endOfData()
	return &r, tok.err()
}

func (FAXCommand) Sentence() string {
	return "#FAX\n"
}

// A FAXResponse is a FAX response.
type FAXResponse struct {
	Date time.Time
}

func parseFAXResponse(tok *tokenizer) (*FAXResponse, error) {
	var r FAXResponse
	year := tok.int()
	month := tok.commaInt()
	day := tok.commaInt()
	r.Date = time.Date(year+1900, time.Month(month)+time.January, day, 0, 0, 0, 0, time.UTC)
	tok.endOfData()
	return &r, tok.err()
}
