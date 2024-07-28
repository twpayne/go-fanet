package fanet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"golang.org/x/text/encoding/charmap"
)

// An FNFResponse is an FNF response.
type FNFResponse struct {
	Source    ID
	Broadcast bool
	Signature Optional[int]
	Type      int
	Payload   []byte
}

// An FNFAck is an FNF ack payload.
type FNFAck struct{}

// An FNFTracking is an FNF tracking payload.
type FNFTracking struct {
	Latitude       float64
	Longitude      float64
	OnlineTracking bool
	Type           int
	Altitude       int
	Speed          float64
	Climb          float64
	Heading        float64
	TurnRate       Optional[float64]
	QNEOffset      Optional[int]
}

// An FNFMessage is an FNF message payload.
type FNFMessage struct {
	SubType uint8
	Message string
}

// An FNFName is an FNF name payload.
type FNFName struct {
	Name string
}

// An FNFGroundTracking is an FNF ground tracking payload.
type FNFGroundTracking struct {
	Latitude       float64
	Longitude      float64
	Type           int
	OnlineTracking bool
}

// An FNFHardwareInfo is an FNF hardware info payload.
type FNFHardwareInfo struct {
	IsPingPongRequest bool
	HardwareSubtype   Optional[int]
	ReleaseVersion    Optional[bool]
	BuildDate         Optional[time.Time]
	ICAOAddress       Optional[[3]byte]
	Uptime            Optional[time.Duration]
	RSSI              Optional[int]
	FANETAddress      Optional[[3]byte]
}

var iso8859_1Decoder = charmap.ISO8859_1.NewDecoder()

func parseFNFResponse(tok *tokenizer) (*FNFResponse, error) {
	var r FNFResponse
	r.Source.Manufacturer = tok.hex()
	r.Source.Device = tok.commaHex()
	r.Broadcast = tok.commaBool()
	r.Signature = tok.commaOptionalHex()
	r.Type = tok.commaHex()
	payloadLength := tok.commaHex()
	r.Payload = tok.commaHexBytes()
	if len(r.Payload) != payloadLength {
		// FIXME improve error message
		// FIXME fix error position, currently returns end of line
		return nil, tok.errOr(fmt.Errorf("payload length mismatch, want %d, got %d", payloadLength, len(r.Payload)))
	}
	tok.endOfData()
	return &r, tok.err()
}

func (r *FNFResponse) ParsePayload() (any, error) {
	switch r.Type {
	case 0:
		// FIXME bounds check
		var fnfAck FNFAck
		return &fnfAck, nil
	case 1:
		// FIXME bounds check
		var fnfTracking FNFTracking
		fnfTracking.Latitude = parseAbsoluteLatitude(r.Payload[:3])
		fnfTracking.Longitude = parseAbsoluteLongitude(r.Payload[3:6])
		onlineTrackingTypeAltitude := binary.LittleEndian.Uint16(r.Payload[6:8])
		fnfTracking.OnlineTracking = onlineTrackingTypeAltitude&0x8000 != 0
		fnfTracking.Type = int(onlineTrackingTypeAltitude & 0x3000 >> 12)
		fnfTracking.Altitude = int(onlineTrackingTypeAltitude & 0x7ff)
		if onlineTrackingTypeAltitude&0x800 != 0 {
			fnfTracking.Altitude <<= 2
		}
		fnfTracking.Speed = float64(r.Payload[8]) / 2
		if r.Payload[8]&0x80 != 0 {
			fnfTracking.Speed *= 5
		}
		fnfTracking.Climb = float64(int8(r.Payload[9]<<1)>>1) / 10
		if r.Payload[9]&0x80 != 0 {
			fnfTracking.Climb *= 5
		}
		fnfTracking.Heading = 360 * float64(r.Payload[10]) / 256
		if len(r.Payload) >= 12 {
			turnRate := float64(int8(r.Payload[11]<<1)>>1) / 4
			if r.Payload[11]&0x80 != 0 {
				turnRate *= 4
			}
			fnfTracking.TurnRate = NewOptional(turnRate)
		}
		if len(r.Payload) >= 13 {
			qneOffset := int(int8(r.Payload[12]<<1) >> 1)
			if r.Payload[12]&0x80 != 0 {
				qneOffset *= 4
			}
			fnfTracking.QNEOffset = NewOptional(qneOffset)
		}
		return &fnfTracking, nil
	case 2:
		var fnfName FNFName
		fnfName.Name = string(r.Payload)
		return &fnfName, nil
	case 3:
		// FIXME bounds check
		var fnfMessage FNFMessage
		fnfMessage.SubType = r.Payload[0]
		fnfMessage.Message = string(r.Payload[1:])
		return &fnfMessage, nil
	case 4:
		// Service
		return nil, errors.ErrUnsupported
	case 5:
		// Landmarks
		return nil, errors.ErrUnsupported
	case 7:
		// FIXME bounds check
		var fnfGroundTracking FNFGroundTracking
		fnfGroundTracking.Latitude = parseAbsoluteLatitude(r.Payload[:3])
		fnfGroundTracking.Longitude = parseAbsoluteLongitude(r.Payload[3:6])
		fnfGroundTracking.Type = int(r.Payload[6] >> 4)
		fnfGroundTracking.OnlineTracking = r.Payload[6]&0x01 != 0
		return &fnfGroundTracking, nil
	case 8:
		// Hardware info, deprecated
		return nil, errors.ErrUnsupported
	case 9:
		// Thermal
		return nil, errors.ErrUnsupported
	case 0x0A:
		if len(r.Payload) < 1 {
			return nil, errors.New("payload too short")
		}
		var fnfHardwareInfo FNFHardwareInfo
		fnfHardwareInfo.IsPingPongRequest = r.Payload[0]&0x80 != 0
		hasHardwareSubTypeAndBuildDate := r.Payload[0]&0x40 != 0
		hasICAOAddress := r.Payload[0]&0x20 != 0
		hasUptime := r.Payload[0]&0x10 != 0
		hasRxRSSI := r.Payload[0]&0x08 != 0
		hasExtendedHeader := r.Payload[0]&0x01 != 0
		minLen := 1
		payloadIndex := 1
		if hasHardwareSubTypeAndBuildDate {
			minLen += 3
		}
		if hasICAOAddress {
			minLen += 3
		}
		if hasUptime {
			minLen += 2
		}
		if hasRxRSSI {
			minLen += 4
		}
		if hasExtendedHeader {
			minLen++
			payloadIndex++
		}
		if len(r.Payload) < minLen {
			return nil, fmt.Errorf("payload is %d bytes, want at least %d", len(r.Payload), minLen)
		}
		if hasHardwareSubTypeAndBuildDate {
			fnfHardwareInfo.HardwareSubtype = NewOptional(int(r.Payload[payloadIndex]))
			year := 2019 + int(r.Payload[payloadIndex+2]>>1&0x3f)
			month := time.Month(r.Payload[payloadIndex+1]>>5 + r.Payload[payloadIndex+2]&0x1<<3)
			day := int(r.Payload[payloadIndex+1] & 0x1f)
			fnfHardwareInfo.BuildDate = NewOptional(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
			fnfHardwareInfo.ReleaseVersion = NewOptional(r.Payload[payloadIndex+2]&0x80 == 0)
			payloadIndex += 3
		}
		if hasICAOAddress {
			var icaoAddress [3]byte
			// FIXME check byte order
			copy(icaoAddress[:], r.Payload[payloadIndex:payloadIndex+3])
			fnfHardwareInfo.ICAOAddress = NewOptional(icaoAddress)
			payloadIndex += 3
		}
		if hasUptime {
			uptimeMinutes := int(r.Payload[payloadIndex]) + int(r.Payload[payloadIndex+1])<<8
			fnfHardwareInfo.Uptime = NewOptional(time.Duration(uptimeMinutes) * time.Minute)
			payloadIndex += 2
		}
		if hasRxRSSI {
			fnfHardwareInfo.RSSI = NewOptional(int(r.Payload[payloadIndex]))
			var fanetAddress [3]byte
			// FIXME check byte order
			copy(fanetAddress[:], r.Payload[payloadIndex+1:payloadIndex+4])
			fnfHardwareInfo.FANETAddress = NewOptional(fanetAddress)
			payloadIndex += 4 //nolint:ineffassign,wastedassign
		}
		return &fnfHardwareInfo, nil
	default:
		return nil, fmt.Errorf("%d: unknown type", r.Type)
	}
}

func (n *FNFName) Payload() (string, error) {
	return iso8859_1Decoder.String(n.Name)
}

func (n *FNFName) FNFType() int {
	return 2
}

func parseAbsoluteLatitude(data []byte) float64 {
	latitude := binary.LittleEndian.Uint32([]byte{data[0], data[1], data[2], 0})
	return float64(int32(latitude<<8)>>8) / 93206
}

func parseAbsoluteLongitude(data []byte) float64 {
	longitude := binary.LittleEndian.Uint32([]byte{data[0], data[1], data[2], 0})
	return float64(int32(longitude<<8)>>8) / 46603
}
