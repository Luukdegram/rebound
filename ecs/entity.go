package ecs

import (
	"sync"
	"sync/atomic"
)

var currID int32 = 0

// Entity is an object that can hold multiple components
// Components can be added/removed/accessed concurrently
type Entity struct {
	id         int32
	m          *sync.RWMutex
	components map[string]Component
	childs     []*Entity
	parent     *Entity
}

//NewEntity creates a new Entity with the given components
func NewEntity(components ...Component) *Entity {
	newID := atomic.AddInt32(&currID, 1)
	comp := make(map[string]Component)
	for _, c := range components {
		if c != nil {
			comp[c.Name()] = c
		}
	}

	return &Entity{
		id:         newID,
		m:          &sync.RWMutex{},
		components: comp,
	}
}

//ID returns the ID of an Entity
func (e *Entity) ID() int32 {
	e.m.RLock()
	val := e.id
	e.m.RUnlock()
	return val
}

// setParent sets the parent of the Entity
func (e *Entity) setParent(p *Entity) {
	e.m.Lock()
	e.parent = p
	e.m.Unlock()
}

// AddChild adds a child to the entity and sets the parent of the child to the current entity
func (e *Entity) AddChild(c *Entity) {
	e.m.Lock()
	e.childs = append(e.childs, c)
	e.m.Unlock()
	c.setParent(e)
}

// AddChildren adds multiple children to an entity and sets the parent of each child to the current Entity
func (e *Entity) AddChildren(childs []*Entity) {
	e.m.Lock()
	e.childs = append(e.childs, childs...)
	e.m.Unlock()
	for _, c := range childs {
		c.setParent(e)
	}
}

// Children returns the children of an Entity
func (e *Entity) Children() []*Entity {
	e.m.RLock()
	val := e.childs
	e.m.RUnlock()
	return val
}

// Parent returns the parent of an Entity
func (e *Entity) Parent() *Entity {
	e.m.RLock()
	val := e.parent
	e.m.RUnlock()
	return val
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
//This returns nil if a [Component] does not exist
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
