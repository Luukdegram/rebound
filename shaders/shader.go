package shaders

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

//ShaderProgram is a generic shader program
type ShaderProgram struct {
	ID               uint32
	vertexShaderID   uint32
	fragmentShaderID uint32
	attributes       []string
}

//NewShaderProgram creates a new shader program given the vertex file and fragment file
func NewShaderProgram(vertexFile string, fragmentFile string) (*ShaderProgram, error) {
	s := new(ShaderProgram)
	var err error

	s.vertexShaderID, err = loadShaderFromFile(vertexFile, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	s.fragmentShaderID, err = loadShaderFromFile(fragmentFile, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	s.ID = gl.CreateProgram()
	gl.AttachShader(s.ID, s.vertexShaderID)
	gl.AttachShader(s.ID, s.fragmentShaderID)

	gl.LinkProgram(s.ID)
	gl.ValidateProgram(s.ID)

	gl.DetachShader(s.ID, s.vertexShaderID)
	gl.DetachShader(s.ID, s.fragmentShaderID)

	return s, nil
}

//getUniformLocation returns the location of the uniform given, returning the OpenGL id as an int32
func (sp ShaderProgram) getUniformLocation(name string) int32 {
	return gl.GetUniformLocation(sp.ID, gl.Str(name+"\x00"))
}

//LoadFloat loads a uniform float into the shader
func (sp ShaderProgram) LoadFloat(name string, value float32) {
	loc := sp.getUniformLocation(name)
	gl.Uniform1f(loc, value)
}

//LoadVec3 loads a uniform Vector into the shader
func (sp ShaderProgram) LoadVec3(name string, value mgl32.Vec3) {
	loc := sp.getUniformLocation(name)
	gl.Uniform3f(loc, value[0], value[1], value[2])
}

//LoadBool loads a boolean into the shader
func (sp ShaderProgram) LoadBool(name string, value bool) {
	loc := sp.getUniformLocation(name)
	var float float32
	if value {
		float = 1
	}
	gl.Uniform1f(loc, float)
}

//LoadMat loads a matrix into the shader
func (sp ShaderProgram) LoadMat(name string, value mgl32.Mat4) {
	loc := sp.getUniformLocation(name)
	gl.UniformMatrix4fv(loc, 1, false, &value[0])
}

//Start starts the shader program
func (sp ShaderProgram) Start() {
	gl.UseProgram(sp.ID)
}

//Stop stops the shader program
func (sp ShaderProgram) Stop() {
	gl.UseProgram(0)
}

//CleanUp deletes the program
func (sp ShaderProgram) CleanUp() {
	gl.DeleteProgram(sp.ID)
}

//LoadShader loads a shader file from system and compiles it
func loadShaderFromFile(fileName string, shaderType uint32) (uint32, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	source, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err
	}

	id, err := compileShader(string(source)+"\x00", shaderType)
	return id, err
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
