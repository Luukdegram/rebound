package rebound

const (
	//PositionComponentName is the name of the PositionComponent
	PositionComponentName = "PositionComponent"
	//RotationComponentName is the name of the RotationComponent
	RotationComponentName = "RotationComponent"
)

//PositionComponent is a 3-element array of float32's
type PositionComponent struct {
	Coords [3]float32
}

// Name returns the name of the PositionComponent
func (p *PositionComponent) Name() string {
	return PositionComponentName
}

// RotationComponent is a float32 array consisting of 3 elements
type RotationComponent struct {
	Coords [3]float32
}

// Name returns the name of the RotationComponent
func (r *RotationComponent) Name() string {
	return RotationComponentName
}
