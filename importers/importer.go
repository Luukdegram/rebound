package importers

import (
	"github.com/luukdegram/rebound/ecs"
)

//Importer is an interface that implements importing files to 3d models
type Importer interface {
	Import(string) (*ecs.Entity, error)
}
