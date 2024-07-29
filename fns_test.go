package fanet

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestFNSRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		name string
		s    string
	}{
		{
			name: "base",
			s:    "#FNS 47.200582,8.523609,200,0.17,0,96,122,11,5,14,1,59\n",
		},
		{
			name: "geoid_separation",
			s:    "#FNS 47.200582,8.523609,200,0.17,0,96,122,11,5,14,1,59,47\n",
		},
		{
			name: "geoid_separation_and_turn_rate",
			s:    "#FNS 47.200582,8.523609,200,0.17,0,96,122,11,5,14,1,59,47,-3\n",
		},
		{
			name: "geoid_separation_and_turn_rate_and_qne_offset",
			s:    "#FNS 47.200582,8.523609,200,0.17,0,96,122,11,5,14,1,59,47,-3,-32\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tok := newTokenizer([]byte(tc.s))
			assert.Equal(t, "#FNS", tok.header())
			command, err := parseFNSCommand(tok)
			assert.NoError(t, err)
			assert.Equal(t, tc.s, command.Sentence())
		})
	}
}
