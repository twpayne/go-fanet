package fanet

import "fmt"

type Response any

type responseParser func(*tokenizer) (any, error)

func makeResponseParser[C any](f func(*tokenizer) (C, error)) responseParser {
	return func(tok *tokenizer) (any, error) {
		return f(tok)
	}
}

var responseParsers = map[string]responseParser{
	"#DBR": makeResponseParser(parseDBRResponse),
	"#DGR": makeResponseParser(parseDGRResponse),
	"#DGV": makeResponseParser(parseDGVResponse),
	"#FAO": makeResponseParser(parseFAOResponse),
	"#FAR": makeResponseParser(parseFARResponse),
	"#FAX": makeResponseParser(parseFAXResponse),
	"#FNA": makeResponseParser(parseFNAResponse),
	"#FNF": makeResponseParser(parseFNFResponse),
	"#FNR": makeResponseParser(parseFNRResponse),
	"#FNZ": makeResponseParser(parseFNZResponse),
}

// ParseResponse parses a response from data.
func ParseResponse(data []byte) (Response, error) {
	tok := newTokenizer(data)
	header := tok.header()
	responseParser, ok := responseParsers[header]
	if !ok {
		return nil, fmt.Errorf("%s: unknown response", header)
	}
	return responseParser(tok)
}

// ParseResponseString parses a response from s.
func ParseResponseString(s string) (Response, error) {
	return ParseResponse([]byte(s))
}
