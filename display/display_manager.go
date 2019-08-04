package display

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

//Manager is an interface that can be implemented to handle the window using i.e. GLFW
type Manager interface {
	Init(width int, height int, title string) error
	ShouldClose() bool
	Close()
	Update()
	GetSize() Size
	RegisterKeyboardHandler(key Key, callback KeyCallback)
	RegisterScrollWheelHandler(callback ScrollCallback)
}

//Size holds width and height of a window
type Size struct {
	Width  int
	Height int
}

//Key is type to hold a keyboard key regardless what window manager is used
type Key int

var glfwKeyToRKey = map[glfw.Key]Key{
	glfw.KeyP: KeyP,
}

//Provide users with a generic set of Keys that will implement the keys of the selected display manager
const (
	Key0 Key = iota
	KeyP
)

//Default returns the default window manager. In this case GLFW
func Default() Manager {
	return NewGLFWManager()
}

//GLFWManager handles the GLFW window
type GLFWManager struct {
	Manager
	w    *glfw.Window
	size *Size
}

//NewGLFWManager creates a new GLFWManager struct
func NewGLFWManager() Manager {
	return &GLFWManager{}
}

//Init initializes the GLFW window
func (g *GLFWManager) Init(width int, height int, title string) error {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		return err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	g.w = window
	g.size = &Size{Width: width, Height: height}

	// initOpenGL initializes OpenGL
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	return nil
}

//KeyCallback is a function that can be triggered when a key is pressed
type KeyCallback func()

//ScrollCallback is a function that can be triggered when the scrollwheel is used
type ScrollCallback func(x, y float64)

//RegisterKeyboardHandler registers a callback action to a certain key
func (g *GLFWManager) RegisterKeyboardHandler(key Key, callback KeyCallback) {
	fn := func(window *glfw.Window, glfwKey glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfwKeyToRKey[glfwKey] && action == glfw.Press {
			callback()
		}
	}
	g.w.SetKeyCallback(fn)
}

//RegisterScrollWheelHandler registers a callback action that triggers when the scrollwheel turns
func (g *GLFWManager) RegisterScrollWheelHandler(callback ScrollCallback) {
	fn := func(window *glfw.Window, x, y float64) {
		callback(x, y)
	}
	g.w.SetScrollCallback(fn)
}

//ShouldClose returns a boolean wether the window should close or not.
// i.e. when the user closes the window.
func (g *GLFWManager) ShouldClose() bool {
	return g.w.ShouldClose()
}

//Close closes the window and terminates GLFW
func (g *GLFWManager) Close() {
	glfw.Terminate()
}

//Update updates the current screen
func (g *GLFWManager) Update() {
	glfw.PollEvents()
	g.w.SwapBuffers()
}

//GetSize returns the size of the window
func (g *GLFWManager) GetSize() Size {
	return *g.size
}
