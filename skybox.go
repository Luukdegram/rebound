package rebound

import "github.com/luukdegram/rebound/geometry"

// Skybox is a cubemap texture that is displayed around the scene
type Skybox struct {
	shader  Shader
	texture uint32
	mesh    *Mesh
}

// NewSkybox creates a new skybox, using the provided faces.
// Expects 6 filenames to be loaded.
// Returns an error if file could not be loaded.
func NewSkybox(files [6]string) (*Skybox, error) {
	tex, err := LoadCubeMap(files)
	if err != nil {
		return nil, err
	}

	s, err := newSkyboxShader()
	if err != nil {
		return nil, err
	}

	m := &Mesh{
		Attributes: []Attribute{Attribute{
			Type: POSITION,
			Size: 3,
			Data: geometry.NewCube(),
		}},
	}

	LoadMesh(m)

	return &Skybox{s, tex, m}, nil
}
