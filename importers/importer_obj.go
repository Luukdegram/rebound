package importers

import (
	"github.com/luukdegram/rebound"
)

//ObjLoader is a simple struct to load and handle .obj files
type ObjLoader struct {
}

//Import imports a .obj file and converts it into a raw model
func (loader *ObjLoader) Import(fileName string) (*rebound.SceneComponent, error) {
	return nil, nil
}
