package ecs

import (
	"sync"
)

var manager *Manager

//Manager controls the [System]'s and updates each one of them, based on their priority
//[System]'s can be added/read/removed concurrently
type Manager struct {
	m       *sync.RWMutex
	systems []System
}

//GetManager returns a [Manager].
func GetManager() *Manager {
	if manager == nil {
		manager = &Manager{
			m:       &sync.RWMutex{},
			systems: []System{},
		}
	}

	return manager
}

//Update handles the update of each [System] based on their priority
func (m *Manager) Update(dt float64) {
	for _, s := range m.systems {
		go s.Update(dt)
	}
}

//AddSystems adds one or more [System]s to the manager
func (m *Manager) AddSystems(systems ...System) {
	m.m.Lock()
	for _, s := range systems {
		m.systems = append(m.systems, s)
	}
	m.m.Unlock()
}

//Systems returns a list of the [System]s the manager contains
func (m *Manager) Systems() []System {
	m.m.RLock()
	val := m.systems
	m.m.RUnlock()
	return val
}

//Swap changes the priority of 2 [System]s by swapping them based on index in the list
func (m *Manager) Swap(i, j int) {
	m.m.Lock()
	m.systems[i], m.systems[j] = m.systems[j], m.systems[i]
	m.m.Unlock()
}
