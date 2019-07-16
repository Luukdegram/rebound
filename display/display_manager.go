package display

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

//Manager is an interface that can be implemented to handle the window using i.e. GLFW
type Manager interface {
	Init(width int, height int, title string) error
	ShouldClose() bool
	Close()
	Update()
}

//Default returns the default window manager. In this case GLFW
func Default() *GLFWManager {
	return NewGLFWManager()
}

//GLFWManager handles the GLFW window
type GLFWManager struct {
	Manager
	w *glfw.Window
}

//NewGLFWManager creates a new GLFWManager struct
func NewGLFWManager() *GLFWManager {
	return new(GLFWManager)
}

//Init initializes the GLFW window
func (g *GLFWManager) Init(width int, height int, title string) error {
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
	return nil
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
