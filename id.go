package fanet

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type ID struct {
	Manufacturer int
	Device       int
}

var (
	errInvalidID = errors.New("invalid ID")

	idRx = regexp.MustCompile(`\A([0-9A-F]{2}):([0-9A-F]{4})\z`)
)

func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id ID) Int() int {
	return id.Manufacturer<<16 + id.Device
}

func (id ID) IsZero() bool {
	return id.Manufacturer == 0 && id.Device == 0
}

func (id *ID) Set(s string) error {
	m := idRx.FindStringSubmatch(s)
	if m == nil {
		return errInvalidID
	}
	manufacturer, _ := strconv.ParseUint(m[1], 16, 8)
	device, _ := strconv.ParseUint(m[2], 16, 16)
	id.Manufacturer = int(manufacturer)
	id.Device = int(device)
	return nil
}

func (id ID) String() string {
	return fmt.Sprintf("%02X:%04X", id.Manufacturer, id.Device)
}

func (ID) Type() string {
	return "id"
}
