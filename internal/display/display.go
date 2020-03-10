package display

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/luukdegram/rebound/input"
	"github.com/luukdegram/rebound/internal/thread"
)

//Display is an interface that can be implemented to handle the window using i.e. GLFW
type Display interface {
	Init(width int, height int, title string) error
	ShouldClose() bool
	Close()
	Update()
	GetSize() Size
}

//Size holds width and height of a window
type Size struct {
	Width  int
	Height int
}

//Default returns the default window manager. In this case GLFW
func Default() Display {
	return NewGLFWDisplay()
}

//GLFWDisplay handles the GLFW window
type GLFWDisplay struct {
	w    *glfw.Window
	size Size
}

//NewGLFWDisplay creates a new GLFWManager struct
func NewGLFWDisplay() Display {
	return &GLFWDisplay{}
}

//Init initializes the GLFW window
func (g *GLFWDisplay) Init(width int, height int, title string) error {
	return thread.CallErr(func() error {
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
			return err
		}

		window.MakeContextCurrent()
		g.w = window
		g.size = Size{Width: width, Height: height}

		// initOpenGL initializes OpenGL
		if err := gl.Init(); err != nil {
			return err
		}
		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)

		// Set key handler
		//g.registerKeyboardHandler()

		return nil
	})
}

var glfwKeyToKey = map[glfw.Key]input.Key{
	glfw.Key0: input.Key0,
	glfw.Key1: input.Key1,
	glfw.Key2: input.Key2,
	glfw.Key3: input.Key3,
	glfw.Key4: input.Key4,
	glfw.Key5: input.Key5,
	glfw.Key6: input.Key6,
	glfw.Key7: input.Key7,
	glfw.Key8: input.Key8,
	glfw.Key9: input.Key9,
	glfw.KeyQ: input.KeyQ,
	glfw.KeyW: input.KeyW,
	glfw.KeyE: input.KeyE,
	glfw.KeyR: input.KeyR,
	glfw.KeyT: input.KeyT,
	glfw.KeyY: input.KeyY,
	glfw.KeyU: input.KeyU,
	glfw.KeyI: input.KeyI,
	glfw.KeyO: input.KeyO,
	glfw.KeyP: input.KeyP,
	glfw.KeyA: input.KeyA,
	glfw.KeyS: input.KeyS,
	glfw.KeyD: input.KeyD,
	glfw.KeyF: input.KeyF,
	glfw.KeyG: input.KeyG,
	glfw.KeyH: input.KeyH,
	glfw.KeyJ: input.KeyJ,
	glfw.KeyK: input.KeyK,
	glfw.KeyL: input.KeyL,
	glfw.KeyZ: input.KeyZ,
	glfw.KeyX: input.KeyX,
	glfw.KeyC: input.KeyC,
	glfw.KeyV: input.KeyV,
	glfw.KeyB: input.KeyB,
	glfw.KeyN: input.KeyN,
	glfw.KeyM: input.KeyM,
}

var glfwMouseButtonToMouseButton = map[glfw.MouseButton]input.MouseButton{
	glfw.MouseButton1: input.LeftMouseButton,
	glfw.MouseButton2: input.RightMouseButton,
	glfw.MouseButton3: input.MiddleMouseButton,
}

//ScrollCallback is a function that can be triggered when the scrollwheel is used
type ScrollCallback func(x, y float64)

//RegisterKeyboardHandler registers a callback action to a certain key
func (g *GLFWDisplay) registerKeyboardHandler() {
	fn := func(window *glfw.Window, glfwKey glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			input.Keys.Set(glfwKeyToKey[glfwKey], true)
		}

		if action == glfw.Release {
			input.Keys.Set(glfwKeyToKey[glfwKey], false)
		}
	}
	g.w.SetKeyCallback(fn)
}

func (g *GLFWDisplay) registerMouseButtonHandler() {
	fn := func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Press {
			input.MouseButtons.Set(glfwMouseButtonToMouseButton[button], true)
		}
		if action == glfw.Release {
			input.MouseButtons.Set(glfwMouseButtonToMouseButton[button], false)
		}
	}
	g.w.SetMouseButtonCallback(fn)
}

//RegisterScrollWheelHandler registers a callback action that triggers when the scrollwheel turns
func (g *GLFWDisplay) RegisterScrollWheelHandler(callback ScrollCallback) {
	fn := func(window *glfw.Window, x, y float64) {
		callback(x, y)
	}
	g.w.SetScrollCallback(fn)
}

//ShouldClose returns a boolean wether the window should close or not.
// i.e. when the user closes the window.
func (g *GLFWDisplay) ShouldClose() bool {
	return g.w.ShouldClose()
}

//Close closes the window and terminates GLFW
func (g *GLFWDisplay) Close() {
	thread.Call(func() {
		glfw.Terminate()
	})
}

//Update updates the current screen
func (g *GLFWDisplay) Update() {
	thread.Call(func() {
		glfw.PollEvents()
		g.w.SwapBuffers()
	})
}

//GetSize returns the size of the window
func (g *GLFWDisplay) GetSize() Size {
	return g.size
}
