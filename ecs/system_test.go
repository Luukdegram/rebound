package ecs

import "testing"

type TestSystem struct {
	BaseSystem
}

func (t *TestSystem) Update(dt float32) {}

type TestComponent struct{}

var (
	componentName                = "TestComponent"
	testComponent *TestComponent = &TestComponent{}
	testEntity    *Entity        = NewEntity(testComponent)
)

func (t *TestComponent) Name() string {
	return componentName
}

type TestChecker struct{}

func (t *TestChecker) Check(e *Entity) bool {
	return e.HasComponent(componentName)
}

func TestAddEntity(t *testing.T) {
	var expect, result int = 1, 0

	ts := &TestSystem{
		NewBaseSystem(),
	}
	ts.AddEntities(&TestChecker{}, testEntity)

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
	ts.AddEntities(&TestChecker{}, testEntity)

	ts.RemoveEntity(testEntity)

	result = len(ts.Entities())

	if result != expect {
		t.Errorf("RemoveEntity failed. Expected %v, but got %v", expect, result)
	}
}
