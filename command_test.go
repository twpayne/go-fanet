package fanet_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-fanet"
)

func TestCommands(t *testing.T) {
	for i, tc := range []struct {
		s           string
		expected    fanet.Command
		expectedErr string
	}{
		{
			s: "#DGJ BLstm\n",
			expected: &fanet.DGJCommand{
				Bootloader: "BLstm",
			},
		},
		{
			s: "#DGL 868,2\n",
			expected: &fanet.DGLCommand{
				Frequency: 868,
				DBm:       2,
			},
		},
		{
			s: "#DGP 1\n",
			expected: &fanet.DGPCommand{
				PowerMode: true,
			},
		},
		{
			s:        "#DGV\n",
			expected: &fanet.DGVCommand{},
		},
		{
			s: "#FAP 1\n",
			expected: &fanet.FAPCommand{
				PowerMode: true,
			},
		},
		{
			s:        "#FNA\n",
			expected: &fanet.FNACommand{},
		},
		{
			s: "#FNC 2,1,3\n",
			expected: &fanet.FNCCommand{
				AircraftType: 2,
				LiveTracking: true,
				GroundType:   3,
			},
		},
		{
			s: "#FNM 1\n",
			expected: &fanet.FNMCommand{
				Mode: 1,
			},
		},
		{
			s: "#FNS 47.200589,8.523513,200,0.17,0,95,122,11,5,14,1,47\n",
			expected: &fanet.FNSCommand{
				Latitude:  47.200589,
				Longitude: 8.523513,
				Altitude:  200,
				SpeedKPH:  0.17,
				Climb:     0,
				Heading:   95,
				Time:      time.Date(2022, time.December, 5, 14, 1, 47, 0, time.UTC),
			},
		},
		{
			s: "#FNS 45.1234,10.5678,500,36,-1.5,45,117,9,30,12,35,2\n",
			expected: &fanet.FNSCommand{
				Latitude:  45.1234,
				Longitude: 10.5678,
				Altitude:  500,
				SpeedKPH:  36,
				Climb:     -1.5,
				Heading:   45,
				Time:      time.Date(2017, time.October, 30, 12, 35, 2, 0, time.UTC),
			},
		},
		{
			s: "#FNS 47.200582,8.523609,200,0.17,0,96\n",
			expected: &fanet.FNSCommand{
				Latitude:  47.200582,
				Longitude: 8.523609,
				Altitude:  200,
				SpeedKPH:  0.17,
				Climb:     0,
				Heading:   96,
			},
		},
		{
			s: "#FNS 47.200582,8.523609,200,0.17,0,96,122,11,5,14,1,59\n",
			expected: &fanet.FNSCommand{
				Latitude:  47.200582,
				Longitude: 8.523609,
				Altitude:  200,
				SpeedKPH:  0.17,
				Climb:     0,
				Heading:   96,
				Time:      time.Date(2022, time.December, 5, 14, 1, 59, 0, time.UTC),
			},
		},
		{
			s: "#FNS 47.200582,8.523609,200,0.17,0,96,122,11,5,14,1,59,50\n",
			expected: &fanet.FNSCommand{
				Latitude:        47.200582,
				Longitude:       8.523609,
				Altitude:        200,
				SpeedKPH:        0.17,
				Climb:           0,
				Heading:         96,
				Time:            time.Date(2022, time.December, 5, 14, 1, 59, 0, time.UTC),
				GeoidSeparation: fanet.NewOptional(50.0),
			},
		},
		{
			s: "#FNS 47.200582,8.523609,200,0.17,0,96,122,11,5,14,1,59,50,15\n",
			expected: &fanet.FNSCommand{
				Latitude:        47.200582,
				Longitude:       8.523609,
				Altitude:        200,
				SpeedKPH:        0.17,
				Climb:           0,
				Heading:         96,
				Time:            time.Date(2022, time.December, 5, 14, 1, 59, 0, time.UTC),
				GeoidSeparation: fanet.NewOptional(50.0),
				TurnRate:        fanet.NewOptional(15.0),
			},
		},
		{
			s: "#FNT 2,0,0,0,0,9,546F6D205061796E65\n",
			expected: &fanet.FNTCommand{
				Type:    2,
				Payload: []byte("Tom Payne"),
			},
		},
		{
			s: "#FNT A,0,0,0,0,6,5012CB0A6700\n",
			expected: &fanet.FNTCommand{
				Type:    0xA,
				Payload: []byte{0x50, 0x12, 0xCB, 0x0A, 0x67, 0x00},
			},
		},
		{
			s:           "#XXX XXX\n",
			expectedErr: "#XXX: unknown command",
		},

		/*
			{
				s: "#DGV build-201709261354\n",
				expected: &fanet.DGVResponse{
					Build:    "build",
					DateCode: "201709261354",
				},
			},
			{
				s: "#DGV build-201709261354\n",
				expected: &fanet.DGVResponse{
					Build:    "build",
					DateCode: "201709261354",
				},
			},
			{
				s: "#DGV build-202008031108\n",
				expected: &fanet.DGVResponse{
					Build:    "build",
					DateCode: "202008031108",
				},
			},
			{
				s: "#FAO 0,DF2029,2,8,47.182989,8.521088,429.0,0.0,-0.1,127.0\n",
				expected: &fanet.FAOResponse{
					Source: 0,
					ID: fanet.ID{
						Manufacturer: 0xdf,
						Device:       0x2029,
					},
					IDType:       2,
					AircraftType: 8,
					Latitude:     47.182989,
					Longitude:    8.521088,
					Altitude:     429,
					GroundSpeed:  0,
					ClimbRate:    -0.1,
					Track:        127,
				},
			},
			{
				s: "#FAX 122,2,1\n",
				expected: &fanet.FAXResponse{
					Date: time.Date(2022, time.March, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			{
				s: "#FNA 11,003F\n",
				expected: &fanet.FNAResponse{
					ID: fanet.ID{
						Manufacturer: 0x11,
						Device:       0x3f,
					},
				},
			},
			{
				s: "#FNF 11,1FE3,1,0,7,7,861A43310F0611\n",
				expected: &fanet.FNFResponse{
					Source: fanet.ID{
						Manufacturer: 0x11,
						Device:       0x1fe3,
					},
					Broadcast: true,
					Signature: 0,
					Type:      7,
					Payload:   []byte{0x86, 0x1A, 0x43, 0x31, 0x0F, 0x06, 0x11},
				},
			},
			{
				s: "#FNF 11,D,1,0,7,7,861A43360F0611\n",
				expected: &fanet.FNFResponse{
					Source: fanet.ID{
						Manufacturer: 0x11,
						Device:       0xd,
					},
					Broadcast: true,
					Signature: 0,
					Type:      7,
					Payload:   []byte{0x86, 0x1A, 0x43, 0x36, 0x0F, 0x06, 0x11},
				},
			},
			{
				s: "#FNF 11,1FE3,1,0,2,C,536B79747261787820322E31\n",
				expected: &fanet.FNFResponse{
					Source: fanet.ID{
						Manufacturer: 0x11,
						Device:       0x1fe3,
					},
					Broadcast: true,
					Signature: 0,
					Type:      2,
					Payload:   []byte("Skytraxx 2.1"),
				},
			},
			{
				s: "#FNF 11,D,1,0,2,C,536B79747261787820332E30\n",
				expected: &fanet.FNFResponse{
					Source: fanet.ID{
						Manufacturer: 0x11,
						Device:       0xd,
					},
					Broadcast: true,
					Signature: 0,
					Type:      2,
					Payload:   []byte("Skytraxx 3.0"),
				},
			},
			{
				s: "#FNF 11,1FE3,0,0,3,4,00596573\n",
				expected: &fanet.FNFResponse{
					Source: fanet.ID{
						Manufacturer: 0x11,
						Device:       0x1fe3,
					},
					Broadcast: false,
					Signature: 0,
					Type:      3,
					Payload:   []byte("\x00Yes"),
				},
			},
			{
				s:           "#FNF 11,D,2,0,7,7,861A43360F0611\n",
				expectedErr: "2: invalid broadcast",
			},
			{
				s:           "#FNF 11,D,0,0,7,7,861A43360F06\n",
				expectedErr: "payload length mismatch, want 7, got 6", // FIXME fix position of error
			},
			{
				s: "#FNZ 1,EU868\n",
				expected: &fanet.FNZResponse{
					Zone:     1,
					ZoneName: "EU868",
				},
			},
		*/
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual, err := fanet.ParseCommandString(tc.s)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, actual)
				assert.Equal(t, tc.expected, actual)
				assert.Equal(t, tc.s, actual.Sentence())
			}
		})
	}
}
