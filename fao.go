package fanet

// An FAOResponse is an FAO response.
type FAOResponse struct {
	Source       int
	ID           ID
	IDType       int
	AircraftType int
	Latitude     float64
	Longitude    float64
	Altitude     float64
	GroundSpeed  float64
	ClimbRate    float64
	Track        float64
}

func parseFAOResponse(tok *tokenizer) (*FAOResponse, error) {
	var r FAOResponse
	r.Source = tok.int()
	id := tok.commaHex()
	r.ID.Manufacturer = id >> 16
	r.ID.Device = id & 0xffff
	r.IDType = tok.commaInt()
	r.AircraftType = tok.commaInt()
	r.Latitude = tok.commaFloat()
	r.Longitude = tok.commaFloat()
	r.Altitude = tok.commaFloat()
	r.GroundSpeed = tok.commaFloat()
	r.ClimbRate = tok.commaFloat()
	if !tok.atEndOfData() {
		r.Track = tok.commaFloat()
	}
	tok.endOfData()
	return &r, tok.err()
}

func (r *FAOResponse) Address() string {
	return "FAO"
}
