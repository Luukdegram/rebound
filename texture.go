package rebound

const (
	//TextureComponentName is the name of the TextureComponent
	TextureComponentName = "TextureComponent"
)

//TextureComponent holds all data required to render a texture
type TextureComponent struct {
	id uint32
	// ShineDamper is the amount of shine applied to the object
	ShineDamper float32
	// Reflectivity is the amount of light is reflected by the object
	Reflectivity float32
	// Transparant turns transparancy on or off on the object
	Transparant bool
	// UseFakeLighting applies lighting to the object without a light source.
	UseFakeLighting bool
}

// Name returns the name of the TextureComponent
func (tc *TextureComponent) Name() string {
	return TextureComponentName
}

// NewTextureComponent creates a new TextureComponent with default values
func NewTextureComponent(id uint32) *TextureComponent {
	return &TextureComponent{
		id:              id,
		ShineDamper:     1000,
		Reflectivity:    1,
		Transparant:    false,
		UseFakeLighting: false,
	}
}
