package models

//TexturedModel is a model with a texture attached
type TexturedModel struct {
	Model   *RawModel
	Texture *ModelTexture
}

//NewTexturedModel creates a new Textured Model and returns it
func NewTexturedModel(model *RawModel, texture *ModelTexture) *TexturedModel {
	return &TexturedModel{Model: model, Texture: texture}
}
