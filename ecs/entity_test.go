package ecs

import "testing"

func TestAddComponent(t *testing.T) {
	var expect, result bool = true, false
	te := NewEntity()
	te.AddComponent(testComponent)

	result = te.HasComponent(componentName)

	if result != expect {
		t.Errorf("AddComponent failed. Expected %v, but got %v", expect, result)
	}
}

func TestRemoveComponent(t *testing.T) {
	var expect, result bool = false, true

	te := NewEntity(testComponent)
	te.RemoveComponent(testComponent)

	result = te.HasComponent(componentName)

	if result != expect {
		t.Errorf("RemoveComponent failed. Expected %v, but got %v", expect, result)
	}
}

func TestID(t *testing.T) {
	te := NewEntity()
	te2 := NewEntity()
	expect := currID

	if te.ID() == te2.ID() {
		t.Errorf("ID failed. Expected different ID's but received the same")
	}

	result := te2.ID()

	if result != expect {
		t.Errorf("ID failed. Expected %v, but got %v", expect, result)
	}
}

func TestComponentFunc(t *testing.T) {
	te := NewEntity(testComponent)

	c := te.Component(componentName)

	if c != testComponent {
		t.Errorf("Component failed. Expect testComponent but got different component")
	}
}
