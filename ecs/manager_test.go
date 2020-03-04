package ecs

import (
	"testing"
)

func TestAddSystem(t *testing.T) {
	var expect, result int = 1, 0

	ts := &TestSystem{}

	manager := GetManager()
	manager.AddSystems(ts)

	result = len(manager.Systems())

	if result != expect {
		t.Errorf("AddSystem failed. Expected %d, but got %d", expect, result)
	}
}
