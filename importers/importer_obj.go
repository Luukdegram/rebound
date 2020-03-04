package importers

import "github.com/luukdegram/rebound/ecs"

//ObjLoader is a simple struct to load and handle .obj files
type ObjLoader struct {
}

//Import imports a .obj file and converts it into a raw model
func (loader *ObjLoader) Import(fileName string) (*ecs.Entity, error) {
	return nil, nil
}
