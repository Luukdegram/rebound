package input

import (
	"sync"
)

// Keys holds the state of each key
var Keys keyManager = keyManager{
	mutex: sync.RWMutex{},
	keys:  make(map[Key]bool),
}

// Mouse holds the state of each mouse button
var Mouse mouseManager = mouseManager{
	mutex:   sync.RWMutex{},
	buttons: make(map[MouseButton]bool),
}

// Cursor represents the cursor on the screen.
var Cursor cursorManager = cursorManager{
	mutex: sync.RWMutex{},
}

// Manager handles the state of input
type Manager interface {
	Set(int, bool)
	get(int) bool
}

// keyManager handles keyboard input
type keyManager struct {
	mutex sync.RWMutex
	keys  map[Key]bool
}

// mouseManager handles mouse input
type mouseManager struct {
	mutex   sync.RWMutex
	buttons map[MouseButton]bool
}

type cursorManager struct {
	mutex sync.RWMutex
	x,
	y,
	lastX,
	lastY float32
}

//Key is type to hold a keyboard key regardless what window manager is used
type Key int

// MouseButton is a type that describes which mouse button is held down
type MouseButton int

// Set the value of a key, based on it being pressed or released
func (km *keyManager) Set(key Key, val bool) {
	km.mutex.Lock()
	km.keys[key] = val
	km.mutex.Unlock()
}

func (km *keyManager) get(key Key) bool {
	km.mutex.Lock()
	defer km.mutex.Unlock()
	val, ok := km.keys[key]
	if !ok {
		km.keys[key] = false
		return false
	}
	return val
}

// Set the value of a key, based on it being pressed or released
func (mm *mouseManager) Set(button MouseButton, val bool) {
	mm.mutex.Lock()
	mm.buttons[button] = val
	mm.mutex.Unlock()
}

func (mm *mouseManager) get(button MouseButton) bool {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	val, ok := mm.buttons[button]
	if !ok {
		mm.buttons[button] = false
		return false
	}
	return val
}

// Set sets the cursor position
func (c *cursorManager) Set(x, y float32) {
	c.mutex.Lock()
	c.lastX, c.lastY = c.x, c.y
	c.x, c.y = x, y
	c.mutex.Unlock()
}

// X returns the X position of the cursor
func (c *cursorManager) X() (x float32) {
	c.mutex.RLock()
	x = c.x
	c.mutex.RUnlock()
	return
}

// Y returns the Y position of the cursor
func (c *cursorManager) Y() (y float32) {
	c.mutex.RLock()
	y = c.y
	c.mutex.RUnlock()
	return
}

// Xoffset returns the offset of X compared to last set
func (c *cursorManager) Xoffset() (x float32) {
	c.mutex.RLock()
	x = c.x - c.lastX
	c.mutex.RUnlock()
	return
}

// Xoffset returns the offset of X compared to last set
func (c *cursorManager) Yoffset() (y float32) {
	c.mutex.RLock()
	y = c.y - c.lastY
	c.mutex.RUnlock()
	return
}

// Pos returns the current X and Y value of the cursor's position
func (c *cursorManager) Pos() (x, y float32) {
	c.mutex.RLock()
	x, y = c.x, c.y
	c.mutex.RUnlock()
	return
}

// Up returns true if the key is currently not pressed
func (k Key) Up() bool {
	return !Keys.get(k)
}

// Down returns true if the key is currently being pressed
func (k Key) Down() bool {
	return Keys.get(k)
}

// Up returns true if the mousebutton is currently not pressed
func (mb MouseButton) Up() bool {
	return !Mouse.get(mb)
}

// Down returns false if the mousebutton is currently not pressed
func (mb MouseButton) Down() bool {
	return Mouse.get(mb)
}

//Provide users with a generic set of Keys that will implement the keys of the selected display manager
const (
	Key0 Key = iota
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeyQ
	KeyW
	KeyE
	KeyR
	KeyT
	KeyY
	KeyU
	KeyI
	KeyO
	KeyP
	KeyA
	KeyS
	KeyD
	KeyF
	KeyG
	KeyH
	KeyJ
	KeyK
	KeyL
	KeyZ
	KeyX
	KeyC
	KeyV
	KeyB
	KeyN
	KeyM
)

// Provide users with a generic set of Mouse Buttons that will be implemented by the display manager
const (
	LeftMouseButton MouseButton = iota
	RightMouseButton
	MiddleMouseButton
)
