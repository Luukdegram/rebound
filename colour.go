package rebound

// Colour is a data structure according to RGBA
type Colour struct {
	R float32
	G float32
	B float32
	A float32
}

// ToSlice creates a slice, where the first element represents Red and the last element the Alpha
func (c *Colour) ToSlice() [4]float32 {
	return [4]float32{c.R, c.G, c.B, c.A}
}

// ColourFromSlice creates a new Colour object from a given slice
func ColourFromSlice(data [4]float32) *Colour {
	return &Colour{data[0], data[1], data[2], data[3]}
}
