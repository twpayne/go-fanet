package fanet

import (
	"encoding/hex"
	"strconv"
	"strings"
)

type sentenceBuilder struct {
	builder strings.Builder
}

func newSentenceBuilder(prefix string) *sentenceBuilder {
	var b sentenceBuilder
	b.builder.WriteString(prefix)
	return &b
}

func (b *sentenceBuilder) String() string {
	return b.builder.String()
}

func (b *sentenceBuilder) bool(value bool) {
	var c byte
	if value {
		c = '1'
	} else {
		c = '0'
	}
	b.builder.WriteByte(c)
}

func (b *sentenceBuilder) bytes(bs []byte) {
	b.builder.WriteString(strings.ToUpper(hex.EncodeToString(bs)))
}

func (b *sentenceBuilder) comma() {
	b.builder.WriteByte(',')
}

func (b *sentenceBuilder) commaBool(value bool) {
	b.comma()
	b.bool(value)
}

func (b *sentenceBuilder) commaBytes(bs []byte) {
	b.comma()
	b.bytes(bs)
}

func (b *sentenceBuilder) commaFloat(f float64) {
	b.comma()
	b.float(f)
}

func (b *sentenceBuilder) commaHex(i int) {
	b.comma()
	b.hex(i)
}

func (b *sentenceBuilder) commaInt(i int) {
	b.comma()
	b.int(i)
}

func (b *sentenceBuilder) commaString(s string) {
	b.comma()
	b.string(s)
}

func (b *sentenceBuilder) float(f float64) {
	b.string(strconv.FormatFloat(f, 'f', -1, 64))
}

func (b *sentenceBuilder) hex(i int) {
	b.string(strings.ToUpper(strconv.FormatInt(int64(i), 16)))
}

func (b *sentenceBuilder) int(i int) {
	b.string(strconv.Itoa(i))
}

func (b *sentenceBuilder) newline() {
	b.builder.WriteByte('\n')
}

func (b *sentenceBuilder) string(s string) {
	b.builder.WriteString(s)
}
