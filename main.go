package main

import (
	"go/build"
	"log"
	"os"
	"runtime"

	"github.com/luukdegram/rebound/models"
	"github.com/luukdegram/rebound/shaders"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	height = 480
	width  = 640
)

var (
	triangle = []float32{
		0, 0.5, 0, // top
		-0.5, -0.5, 0, // left
		0.5, -0.5, 0, // right
	}
	indices = []uint32{
		0, 1, 2,
	}
	textureCoords = []float32{
		0, 0,
		0, 1,
		1, 1,
	}
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	runtime.LockOSThread()

	dir, err := importPathToDir("github.com/luukdegram/rebound/assets/")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

func main() {
	window := initGlfw()
	defer glfw.Terminate()

	initOpenGL()
	renderer := new(Renderer)
	model := loadToVAO(triangle, textureCoords, indices)
	texture, err := loadTexture("textures/square.png")
	if err != nil {
		panic(err)
	}
	modelTexture := models.NewModelTexture(texture)
	texturedModel := models.NewTexturedModel(model, modelTexture)

	shader, err := shaders.NewStaticShader("shaders/vertexShader.vert", "shaders/fragmentShader.frag")
	if err != nil {
		panic(err)
	}
	for !window.ShouldClose() {
		renderer.Prepare()
		shader.ShaderProgram.Start()
		renderer.Render(texturedModel)
		shader.ShaderProgram.Stop()

		glfw.PollEvents()
		window.SwapBuffers()
	}
	cleanUp()
}

func initGlfw() *glfw.Window {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	return window
}

// initOpenGL initializes OpenGL
func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}
