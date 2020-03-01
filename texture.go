package rebound

const (
	//TextureComponentName is the name of the TextureComponent
	TextureComponentName = "TextureComponent"
)

//TextureComponent holds all data required to render a texture
type TextureComponent struct {
	id              uint32
	ShineDamper     float32
	Reflectivity    float32
	Transparancy    bool
	UseFakeLighting bool
}

//Name returns the name of the TextureComponent
func (tc *TextureComponent) Name() string {
	return TextureComponentName
}
