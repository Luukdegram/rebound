package ecs

//Component holds data and can be added to an [Entity]
type Component interface {
	//Name returns the name of a [Component]
	Name() string
}
