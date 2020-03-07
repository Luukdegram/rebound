package ecs

import "testing"

type TestSystem struct {
	BaseSystem
}

func (t *TestSystem) Update(dt float32) {}
func (t *TestSystem) Name() string {
	return "TestSystem"
}

type TestComponent struct{}

var (
	componentName                = "TestComponent"
	testComponent *TestComponent = &TestComponent{}
	testEntity    *Entity        = NewEntity(testComponent)
)

func (t *TestComponent) Name() string {
	return componentName
}

func TestAddEntity(t *testing.T) {
	var expect, result int = 1, 0

	ts := &TestSystem{
		NewBaseSystem(),
	}
	ts.AddEntities(testEntity)

	result = len(ts.Entities())

	if result != expect {
		t.Errorf("Failed AddEntity. Expected %d, but got %d.", expect, result)
	}
}

func TestRemoveEntity(t *testing.T) {
	var expect, result int = 0, 1

	ts := &TestSystem{
		NewBaseSystem(),
	}
	ts.AddEntities(testEntity)

	ts.RemoveEntity(testEntity)

	result = len(ts.Entities())

	if result != expect {
		t.Errorf("RemoveEntity failed. Expected %v, but got %v", expect, result)
	}
}
