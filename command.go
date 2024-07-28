package fanet

import "fmt"

type Command interface {
	Sentence() string
}

type commandParser func(*tokenizer) (Command, error)

func makeCommandParser[C Command](f func(*tokenizer) (C, error)) commandParser {
	return func(tok *tokenizer) (Command, error) {
		return f(tok)
	}
}

var commandParsers = map[string]commandParser{
	"#DGJ": makeCommandParser(parseDGJCommand),
	"#DGL": makeCommandParser(parseDGLCommand),
	"#DGP": makeCommandParser(parseDGPCommand),
	"#DGV": makeCommandParser(parseDGVCommand),
	"#FAP": makeCommandParser(parseFAPCommand),
	"#FAT": makeCommandParser(parseFATCommand),
	"#FAX": makeCommandParser(parseFAXCommand),
	"#FNA": makeCommandParser(parseFNACommand),
	"#FNC": makeCommandParser(parseFNCCommand),
	"#FNM": makeCommandParser(parseFNMCommand),
	"#FNS": makeCommandParser(parseFNSCommand),
	"#FNT": makeCommandParser(parseFNTCommand),
}

// ParseCommand parses a response from data.
func ParseCommand(data []byte) (Command, error) {
	tok := newTokenizer(data)
	header := tok.header()
	commandParser, ok := commandParsers[header]
	if !ok {
		return nil, fmt.Errorf("%s: unknown command", header)
	}
	return commandParser(tok)
}

// ParseCommandString parses a response from s.
func ParseCommandString(s string) (Command, error) {
	return ParseCommand([]byte(s))
}
