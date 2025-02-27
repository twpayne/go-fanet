package fanet

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	errExpectedBinaryDigit = errors.New("expected binary digit")
	errExpectedComma       = errors.New("expected comma")
	errExpectedDigit       = errors.New("expected digit")
	errExpectedEndOfData   = errors.New("expected end of data")
	errExpectedFloat       = errors.New("expected float")
	errExpectedHexDigit    = errors.New("expected hex digit")
	errExpectedHeader      = errors.New("expected header")
	errExpectedRegexp      = errors.New("expected regexp")
	errUnexpectedEndOfData = errors.New("unexpected end of data")

	floatRegexp = regexp.MustCompile(`\A-?\d+(?:\.\d*)?`)
)

type SyntaxError struct {
	Data []byte
	Pos  int
	Err  error
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("syntax error at position %d: %v", e.Pos, e.Err)
}

func (e *SyntaxError) Unwrap() error {
	return e.Err
}

type tokenizer struct {
	data  []byte
	pos   int
	error error
}

func newTokenizer(data []byte) *tokenizer {
	return &tokenizer{
		data: data,
	}
}

func (t *tokenizer) atEndOfData() bool {
	if t.error != nil {
		return true
	}
	// FIXME better end of line checking
	return t.pos == len(t.data) || t.pos == len(t.data)-1 && t.data[t.pos] == '\n'
}

func (t *tokenizer) bool() bool {
	if t.error != nil {
		return false
	}
	if t.pos == len(t.data) {
		t.error = errUnexpectedEndOfData
		return false
	}
	switch t.data[t.pos] {
	case '0':
		t.pos++
		return false
	case '1':
		t.pos++
		return true
	default:
		t.error = errExpectedBinaryDigit
		return false
	}
}

func (t *tokenizer) bytes() []byte {
	if t.error != nil {
		return nil
	}
	if t.pos == len(t.data) {
		return nil
	}
	start := t.pos
	for t.pos < len(t.data) && t.data[t.pos] != ',' && t.data[t.pos] != '\n' {
		t.pos++
	}
	return t.data[start:t.pos]
}

func (t *tokenizer) comma() {
	if t.error != nil {
		return
	}
	if t.pos == len(t.data) {
		t.error = errUnexpectedEndOfData
		return
	}
	if t.data[t.pos] != ',' {
		t.error = errExpectedComma
		return
	}
	t.pos++
}

func (t *tokenizer) commaBool() bool {
	t.comma()
	return t.bool()
}

func (t *tokenizer) commaFloat() float64 {
	t.comma()
	return t.float()
}

func (t *tokenizer) commaHex() int {
	t.comma()
	return t.hex()
}

func (t *tokenizer) commaHexBytes() []byte {
	t.comma()
	return t.hexBytes()
}

func (t *tokenizer) commaInt() int {
	t.comma()
	return t.int()
}

func (t *tokenizer) commaOptionalHex() Optional[int] {
	t.comma()
	return t.optionalHex()
}

func (t *tokenizer) commaString() string {
	t.comma()
	return t.string()
}

func (t *tokenizer) endOfData() {
	if t.error != nil {
		return
	}
	if t.pos != len(t.data) && (t.pos != len(t.data)-1 || t.data[t.pos] != '\n') {
		t.error = errExpectedEndOfData
		return
	}
}

func (t *tokenizer) err() error {
	if t.error == nil {
		return nil
	}
	return &SyntaxError{
		Data: t.data,
		Pos:  t.pos,
		Err:  t.error,
	}
}

func (t *tokenizer) errOr(err error) error {
	if t.error != nil {
		return t.err()
	}
	return err
}

func (t *tokenizer) float() float64 {
	if t.error != nil {
		return 0
	}
	if t.pos == len(t.data) {
		t.error = errUnexpectedEndOfData
		return 0
	}
	m := t.regexp(floatRegexp)
	if m == nil {
		t.error = errExpectedFloat
		return 0
	}
	value, _ := strconv.ParseFloat(string(m[0]), 64)
	return value
}

func (t *tokenizer) header() string {
	if t.error != nil {
		return ""
	}
	if t.pos == len(t.data) {
		t.error = errExpectedHeader
		return ""
	}
	start := t.pos
	if t.data[start] != '#' {
		t.error = errExpectedHeader
		return ""
	}
	for i := start + 1; i < len(t.data); i++ {
		if t.data[i] == ' ' || t.data[i] == '\n' {
			t.pos = i + 1
			return string(t.data[start:i])
		}
	}
	t.error = errExpectedHeader
	return ""
}

func (t *tokenizer) hex() int {
	if t.error != nil {
		return 0
	}
	if t.pos == len(t.data) {
		t.error = errUnexpectedEndOfData
		return 0
	}
	value, ok := hexDigitValue(t.data[t.pos])
	if !ok {
		t.error = errExpectedHexDigit
		return 0
	}
	t.pos++
	for t.pos < len(t.data) {
		hexDigit, ok := hexDigitValue(t.data[t.pos])
		if !ok {
			break
		}
		value = 16*value + hexDigit
		t.pos++
	}
	return value
}

func (t *tokenizer) hexBytes() []byte {
	if t.error != nil {
		return nil
	}
	if t.pos == len(t.data) {
		return nil
	}
	value := []byte{}
	for t.pos+1 < len(t.data) {
		hexDigit1, ok := hexDigitValue(t.data[t.pos])
		if !ok {
			t.error = errExpectedHexDigit
			return nil
		}
		t.pos++
		hexDigit2, ok := hexDigitValue(t.data[t.pos])
		if !ok {
			t.error = errExpectedHexDigit
			return nil
		}
		t.pos++
		byteValue := hexDigit1<<4 + hexDigit2
		value = append(value, byte(byteValue))
	}
	return value
}

func (t *tokenizer) int() int {
	if t.error != nil {
		return 0
	}
	if t.pos == len(t.data) {
		t.error = errUnexpectedEndOfData
		return 0
	}
	sign := 1
	if t.data[t.pos] == '-' {
		sign = -1
		t.pos++
	}
	return sign * t.unsignedInt()
}

func (t *tokenizer) optionalHex() Optional[int] {
	if t.error != nil {
		return Optional[int]{}
	}
	if t.pos == len(t.data) {
		t.error = errUnexpectedEndOfData
		return Optional[int]{}
	}
	if t.data[t.pos] == ',' {
		t.pos++
		return Optional[int]{}
	}
	return NewOptional(t.hex())
}

func (t *tokenizer) regexp(regexp *regexp.Regexp) [][]byte {
	if t.error != nil {
		return nil
	}
	m := regexp.FindSubmatch(t.data[t.pos:])
	if m == nil {
		t.error = errExpectedRegexp
		return nil
	}
	t.pos += len(m[0])
	return m
}

func (t *tokenizer) rest() string {
	if t.error != nil {
		return ""
	}
	start := t.pos
	for i := t.pos; i < len(t.data); i++ {
		if t.data[i] == '\n' {
			t.pos = i
			return string(t.data[start:i])
		}
	}
	t.error = errExpectedEndOfData
	result := string(t.data[t.pos:])
	t.pos = len(t.data)
	return result
}

func (t *tokenizer) string() string {
	bytes := t.bytes()
	if bytes == nil {
		return ""
	}
	return string(bytes)
}

func (t *tokenizer) unsignedInt() int {
	if t.error != nil {
		return 0
	}
	if t.pos == len(t.data) {
		t.error = errUnexpectedEndOfData
		return 0
	}
	if t.data[t.pos] < '0' || '9' < t.data[t.pos] {
		t.error = errExpectedDigit
		return 0
	}
	value := int(t.data[t.pos] - '0')
	t.pos++
	for t.pos < len(t.data) {
		digit, ok := digitValue(t.data[t.pos])
		if !ok {
			break
		}
		value = 10*value + digit
		t.pos++
	}
	return value
}

func digitValue(c byte) (int, bool) {
	if '0' <= c && c <= '9' {
		return int(c - '0'), true
	}
	return 0, false
}

func hexDigitValue(c byte) (int, bool) {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0'), true
	case 'A' <= c && c <= 'F':
		return int(c - 'A' + 10), true
	case 'a' <= c && c <= 'f':
		return int(c - 'a' + 10), true
	default:
		return 0, false
	}
}
