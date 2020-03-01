package ecs

import (
	"sync"
	"sync/atomic"
)

var currID int32 = 0

//Entity is an object that can hold multiple components
//[Component]s can be added/removed/accessed concurrently
type Entity struct {
	id         int32
	m          *sync.RWMutex
	components map[string]Component
}

//NewEntity creates a new Entity with the given components
func NewEntity(components ...Component) *Entity {
	newID := atomic.AddInt32(&currID, 1)
	comp := make(map[string]Component)
	for _, c := range components {
		comp[c.Name()] = c
	}

	return &Entity{
		id:         newID,
		m:          &sync.RWMutex{},
		components: comp,
	}
}

//ID returns the ID of an Entity
func (e *Entity) ID() int32 {
	return e.id
}

//AddComponent adds a [Component] to the Entity
func (e *Entity) AddComponent(c Component) {
	e.m.Lock()
	e.components[c.Name()] = c
	e.m.Unlock()
}

//RemoveComponent removes the given [Component] from the Entity
func (e *Entity) RemoveComponent(c Component) {
	e.m.Lock()
	delete(e.components, c.Name())
	e.m.Unlock()
}

//Component returns a [Component] based on its name
//This returns false if a [Component] does not exist
func (e *Entity) Component(name string) Component {
	e.m.RLock()
	val, _ := e.components[name]
	e.m.RUnlock()
	return val
}

//HasComponent checks if a [Component] exists based on the given name
func (e *Entity) HasComponent(name string) bool {
	e.m.RLock()
	_, exists := e.components[name]
	e.m.RUnlock()
	return exists
}
