package ecs

import (
	"sync"
)

//System handles updates on entities based on the components they withold
type System interface {
	//Update is run on every system to update the entities
	Update(dt float32)
	//AddEntities allow for 1 or more entities to be added to the system
	AddEntities(EntityChecker, ...*Entity)
	//Entities returns a map of entities the system contains
	Entities() map[int32]*Entity
	//RemoveEntity removes a singular entity from the system
	RemoveEntity(e *Entity)
}

//EntityChecker is an interface to allow for entity checking to see if they meet the system requirements
type EntityChecker interface {
	//Check returns true if an entity is valid for the current system. I.E. if it contains the needed components
	Check(e *Entity) bool
}

//BaseSystem is a base implementation of [System]. However, requires an [Update] function to meet the requirements of the interface
type BaseSystem struct {
	entities map[int32]*Entity
	m        *sync.RWMutex
}

//NewBaseSystem creates a new [BaseSystem] that allows for concurrent access
func NewBaseSystem() BaseSystem {
	return BaseSystem{
		m:        &sync.RWMutex{},
		entities: make(map[int32]*Entity),
	}
}

//AddEntities adds one or multiple entities to the system using the provided checker.
//It only adds entities that meet the checker's requirements.
func (bs *BaseSystem) AddEntities(ec EntityChecker, entities ...*Entity) {
	for _, e := range entities {
		if ec.Check(e) {
			bs.entities[e.id] = e
		}
	}

}

//Entities returns a map of entities where its key is the [ID()] of the Entity
func (bs *BaseSystem) Entities() map[int32]*Entity {
	bs.m.RLock()
	val := bs.entities
	bs.m.RUnlock()
	return val
}

//RemoveEntity removes the given [Entity] from the [System]
func (bs *BaseSystem) RemoveEntity(e *Entity) {
	bs.m.Lock()
	delete(bs.entities, e.id)
	bs.m.Unlock()
}
