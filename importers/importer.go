package importers

import "github.com/luukdegram/rebound"

//Importer is an interface that implements importing files to 3d models
type Importer interface {
	Import(string) (*rebound.SceneComponent, error)
}
