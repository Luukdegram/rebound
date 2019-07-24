package main

import (
	"go/build"
	"log"
	"os"
	"rebound/display"
	"runtime"

	"rebound/models"
	"rebound/shaders"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	height = 480
	width  = 640
)

var (
	triangle = []float32{
		-0.5, 0.5, 0, // V1
		-0.5, -0.5, 0, // V2
		0.5, -0.5, 0, // V3
		0.5, 0.5, 0, // V4

	}
	indices = []uint32{
		0, 1, 3, //Top left Triangle
		3, 1, 2, //Bottom right Triangle
	}
	textureCoords = []float32{
		0, 0, //V0
		0, 1, //V1
		1, 1, //V2
		1, 0, //V3
	}
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	runtime.LockOSThread()

	dir, err := importPathToDir("rebound/assets/")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

func main() {
	window := display.Manager(display.Default())

	err := window.Init(width, height, "Rebound Engine")
	if err != nil {
		panic(err)
	}
	defer window.Close()

	initOpenGL()
	renderer := new(Renderer)
	model := LoadToVAO(triangle, textureCoords, indices)
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

		window.Update()
	}
	shader.ShaderProgram.CleanUp()
	cleanUp()
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
