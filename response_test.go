package fanet_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-fanet"
)

func TestParseResponse(t *testing.T) {
	for i, tc := range []struct {
		s           string
		expected    fanet.Response
		expectedErr string
	}{
		{
			s: "#DBR OK\n",
			expected: &fanet.DBRResponse{
				Type: "OK",
			},
		},
		{
			s: "#DBR ERR,1,unknown DBG command\n",
			expected: &fanet.DBRResponse{
				Type:      "ERR",
				Status:    1,
				StatusStr: "unknown DBG command",
			},
		},
		{
			s: "#DGR OK\n",
			expected: &fanet.DGRResponse{
				Type: "OK",
			},
		},
		{
			s: "#DGR ERR,70,power switch failed\n",
			expected: &fanet.DGRResponse{
				Type:      "ERR",
				Status:    70,
				StatusStr: "power switch failed",
			},
		},
		{
			s:           "#DGR XXX\n",
			expectedErr: "XXX: unknown DGR type",
		},
		{
			s: "#DGV 1.06, 1722e538e\n",
			expected: &fanet.DGVResponse{
				BuildDateCode: "1.06, 1722e538e",
			},
		},
		{
			s: "#FAR OK\n",
			expected: &fanet.FARResponse{
				Type: "OK",
			},
		},
		{
			s: "#FAR ERR,91,FLARM expired\n",
			expected: &fanet.FARResponse{
				Type:      "ERR",
				Status:    91,
				StatusStr: "FLARM expired",
			},
		},
		{
			s:           "#FAR XXX\n",
			expectedErr: "XXX: unknown FAR type",
		},
		{
			s: "#FAX 125,2,1\n",
			expected: &fanet.FAXResponse{
				Date: time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			s: "#FNA 11,ABCD\n",
			expected: &fanet.FNAResponse{
				ID: fanet.ID{
					Manufacturer: 0x11,
					Device:       0xABCD,
				},
			},
		},
		{
			s: "#FNR OK\n",
			expected: &fanet.FNRResponse{
				Type: "OK",
			},
		},
		{
			s: "#FNR MSG,1,initialized\n",
			expected: &fanet.FNRResponse{
				Type:      "MSG",
				Status:    1,
				StatusStr: "initialized",
			},
		},
		{
			s: "#FNR ERR,12,incompatible type\n",
			expected: &fanet.FNRResponse{
				Type:      "ERR",
				Status:    12,
				StatusStr: "incompatible type",
			},
		},
		{
			s: "#FNR MSG,13,power down\n",
			expected: &fanet.FNRResponse{
				Type:      "MSG",
				Status:    13,
				StatusStr: "power down",
			},
		},
		{
			s: "#FNR ERR,10,no source address\n",
			expected: &fanet.FNRResponse{
				Type:      "ERR",
				Status:    10,
				StatusStr: "no source address",
			},
		},
		{
			s: "#FNR ERR,14,tx buffer full\n",
			expected: &fanet.FNRResponse{
				Type:      "ERR",
				Status:    14,
				StatusStr: "tx buffer full",
			},
		},
		{
			s: "#FNR ERR,30,too short\n",
			expected: &fanet.FNRResponse{
				Type:      "ERR",
				Status:    30,
				StatusStr: "too short",
			},
		},
		{
			s: "#FNR ACK,20,12F2\n",
			expected: &fanet.FNRResponse{
				Type: "ACK",
				Destination: fanet.ID{
					Manufacturer: 0x20,
					Device:       0x12f2,
				},
			},
		},
		{
			s:           "#FNR XXX\n",
			expectedErr: "XXX: unknown FNR type",
		},
		{
			s:           "#XXX XXX\n",
			expectedErr: "#XXX: unknown response",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual, err := fanet.ParseResponseString(tc.s)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, actual)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
