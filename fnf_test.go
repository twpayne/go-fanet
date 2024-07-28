package fanet_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-fanet"
)

func TestFNF_ParsePayload(t *testing.T) {
	for _, tc := range []struct {
		s           string
		expected    any
		expectedErr error
	}{
		{
			s: "#FNF 20,C9E,1,0,1,B,601A43330F06B91100008C\n",
			expected: &fanet.FNFTracking{
				Latitude:       47.18219857090745,
				Longitude:      8.521060875909276,
				OnlineTracking: false,
				Type:           1,
				Altitude:       441,
				Speed:          0,
				Climb:          0,
				Heading:        196.875,
			},
		},
		{
			s: "#FNF 11,D,1,0,2,C,536B79747261787820332E30\n",
			expected: &fanet.FNFName{
				Name: "Skytraxx 3.0",
			},
		},
		{
			s: "#FNF 11,1FE3,1,0,2,C,536B79747261787820322E31\n",
			expected: &fanet.FNFName{
				Name: "Skytraxx 2.1",
			},
		},
		{
			s: "#FNF A,493,1,0,2,9,546F6D205061796E65\n",
			expected: &fanet.FNFName{
				Name: "Tom Payne",
			},
		},
		{
			s:           "#FNF E8,1412,1,0,5,B,C4D7FC5CC5227B9B0C22DC\n",
			expectedErr: errors.ErrUnsupported,
		},
		{
			s: "#FNF 11,1FE3,1,0,7,7,8B1A432B0F0611\n",
			expected: &fanet.FNFGroundTracking{
				Latitude:       47.18265991459777,
				Longitude:      8.520889213140785,
				Type:           1,
				OnlineTracking: true,
			},
		},
		{
			s: "#FNF A,493,1,0,7,7,841A43310F0611\n",
			expected: &fanet.FNFGroundTracking{
				Latitude:       47.18258481213656,
				Longitude:      8.521017960217153,
				Type:           1,
				OnlineTracking: true,
			},
		},
		{
			s:           "#FNF 11,D,1,0,8,5,01DE062014\n",
			expectedErr: errors.ErrUnsupported,
		},
		{
			s: "#FNF A,493,1,0,A,6,5012670A0A00\n",
			expected: &fanet.FNFHardwareInfo{
				HardwareSubtype: fanet.NewOptional(0x12),
				ReleaseVersion:  fanet.NewOptional(true),
				BuildDate:       fanet.NewOptional(time.Date(2024, time.March, 7, 0, 0, 0, 0, time.UTC)),
				Uptime:          fanet.NewOptional(10 * time.Minute),
			},
		},
		{
			s: "#FNF A,493,1,0,A,6,5012680A0B00\n",
			expected: &fanet.FNFHardwareInfo{
				HardwareSubtype: fanet.NewOptional(0x12),
				ReleaseVersion:  fanet.NewOptional(true),
				BuildDate:       fanet.NewOptional(time.Date(2024, time.March, 8, 0, 0, 0, 0, time.UTC)),
				Uptime:          fanet.NewOptional(11 * time.Minute),
			},
		},
	} {
		t.Run(tc.s, func(t *testing.T) {
			fnf, err := fanet.ParseResponseString(tc.s)
			assert.NoError(t, err)
			actual, err := fnf.(*fanet.FNFResponse).ParsePayload()
			if tc.expectedErr != nil {
				assert.EqualError(t, err, "unsupported operation")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
