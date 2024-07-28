package fanet

import (
	"time"
)

// An FNSCommand is an FNS command.
type FNSCommand struct {
	Latitude        float64
	Longitude       float64
	Altitude        float64
	SpeedKPH        float64
	Climb           float64
	Heading         float64
	Time            time.Time
	GeoidSeparation Optional[float64]
	TurnRate        Optional[float64]
}

func parseFNSCommand(tok *tokenizer) (*FNSCommand, error) {
	var c FNSCommand
	c.Latitude = tok.float()
	c.Longitude = tok.commaFloat()
	c.Altitude = tok.commaFloat()
	c.SpeedKPH = tok.commaFloat()
	c.Climb = tok.commaFloat()
	c.Heading = tok.commaFloat()
	if !tok.atEndOfData() {
		year := tok.commaInt()
		month := tok.commaInt()
		day := tok.commaInt()
		hour := tok.commaInt()
		minute := tok.commaInt()
		second := tok.commaInt()
		c.Time = time.Date(year+1900, time.Month(month)+time.January, day, hour, minute, second, 0, time.UTC)
	}
	if !tok.atEndOfData() {
		c.GeoidSeparation = NewOptional(tok.commaFloat())
	}
	if !tok.atEndOfData() {
		c.TurnRate = NewOptional(tok.commaFloat())
	}
	tok.endOfData()
	return &c, tok.err()
}

func (c *FNSCommand) Sentence() string {
	b := newSentenceBuilder("#FNS ")
	b.float(c.Latitude)
	b.commaFloat(c.Longitude)
	b.commaFloat(c.Altitude)
	b.commaFloat(c.SpeedKPH)
	b.commaFloat(c.Climb)
	b.commaFloat(c.Heading)
	if !c.Time.IsZero() {
		b.commaInt(c.Time.Year() - 1900)
		b.commaInt(int(c.Time.Month()) - 1)
		b.commaInt(c.Time.Day())
		b.commaInt(c.Time.Hour())
		b.commaInt(c.Time.Minute())
		b.commaInt(c.Time.Second())
	} else if c.Time.IsZero() && (c.GeoidSeparation.Valid || c.TurnRate.Valid) {
		b.string(",,,,,,")
	}
	if c.GeoidSeparation.Valid {
		b.commaFloat(c.GeoidSeparation.Value)
		if c.TurnRate.Valid {
			b.commaFloat(c.TurnRate.Value)
		}
	} else if c.TurnRate.Valid {
		b.comma()
		b.commaFloat(c.TurnRate.Value)
	}
	b.newline()
	return b.String()
}
