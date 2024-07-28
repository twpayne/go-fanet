// Package fanet generates and parses FANET sentences.
//
// See https://github.com/3s1d/fanet-stm32.
// See https://github.com/3s1d/fanet-stm32/blob/master/fanet_module.pdf.
// See https://github.com/3s1d/fanet-stm32/blob/master/Src/fanet/radio/protocol.txt.
package fanet

var AircraftTypes = map[int]string{
	0: "Other",
	1: "Paraglider",
	2: "Hangglider",
	3: "Balloon",
	4: "Glider",
	5: "Powered Aircraft",
	6: "Helicopter",
	7: "UAV",
}

var ManufacturerNames = map[int]string{
	0x00: "[reserved]",
	0x01: "Skytraxx",
	0x03: "BitBroker.eu",
	0x04: "AirWhere",
	0x05: "Windline",
	0x06: "Burnair.ch",
	0x07: "SoftRF",
	0x08: "GXAircom",
	0x09: "Airtribune",
	0x10: "alfapilot",
	0x0A: "FLARM",
	0x11: "FANET+",
	0x20: "XC Tracer",
	0xCB: "CloudBuddy",
	0xE0: "OGN Tracker",
	0xE4: "4aviation",
	0xFA: "Various",
	0xFB: "Expressif",
	0xFC: "Unregistered",
	0xFD: "Unregistered",
	0xFE: "[Multicast]",
	0xFF: "[reserved]",
}

var GroundTrackingTypes = map[int]string{
	0:  "Other",
	1:  "Walking",
	2:  "Vehicle",
	3:  "Bike",
	4:  "Boot", // Boat
	8:  "Need a ride",
	9:  "Landed well",
	12: "Need technical support",
	13: "Need medical help",
	14: "Distress call",
	15: "Distress call automatically",
}
