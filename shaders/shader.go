package shaders

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

//ShaderComponentName is the name of the ShaderComponent
const ShaderComponentName string = "ShaderComponent"

//ShaderComponent holds the data a ShaderProgram requires to use given shader
type ShaderComponent struct {
	id               uint32
	vertexShaderID   uint32
	fragmentShaderID uint32
}

//Name returns the name of the ShaderComponent
//Needed to fit the Component interface
func (sc *ShaderComponent) Name() string {
	return ShaderComponentName
}

//NewShaderComponent returns a new ShaderComponent by compiling the given vertexShader and fragmentShader
//Returns an error if any of the shaders could not be compiled
func NewShaderComponent(vertexShader, fragmentShader string) (*ShaderComponent, error) {
	var err error
	s := &ShaderComponent{}

	if s.vertexShaderID, err = compileShader(vertexShader+"\x00", gl.VERTEX_SHADER); err != nil {
		return nil, err
	}

	if s.fragmentShaderID, err = compileShader(fragmentShader+"\x00", gl.FRAGMENT_SHADER); err != nil {
		return nil, err
	}

	s.id = gl.CreateProgram()
	gl.AttachShader(s.id, s.vertexShaderID)
	gl.AttachShader(s.id, s.fragmentShaderID)

	gl.LinkProgram(s.id)
	gl.ValidateProgram(s.id)

	gl.DetachShader(s.id, s.vertexShaderID)
	gl.DetachShader(s.id, s.fragmentShaderID)

	return s, nil
}

//GetUniformLocation returns the location of the uniform given, returning the OpenGL id as an int32
func GetUniformLocation(s ShaderComponent, name string) int32 {
	return gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
}

//LoadFloat loads a uniform float into the shader
func LoadFloat(s ShaderComponent, name string, value float32) {
	loc := GetUniformLocation(s, name)
	gl.Uniform1f(loc, value)
}

//LoadVec3 loads a uniform Vector into the shader
func LoadVec3(s ShaderComponent, name string, value [3]float32) {
	loc := GetUniformLocation(s, name)
	gl.Uniform3f(loc, value[0], value[1], value[2])
}

//LoadBool loads a boolean into the shader
func LoadBool(s ShaderComponent, name string, value bool) {
	loc := GetUniformLocation(s, name)
	var float float32
	if value {
		float = 1
	}
	gl.Uniform1f(loc, float)
}

//LoadMat loads a matrix into the shader
func LoadMat(s ShaderComponent, name string, value [16]float32) {
	loc := GetUniformLocation(s, name)
	gl.UniformMatrix4fv(loc, 1, false, &value[0])
}

//Start starts the shader program
func Start(s ShaderComponent) {
	gl.UseProgram(s.id)
}

//Stop stops the shader program
func Stop() {
	gl.UseProgram(0)
}

//CleanUp deletes the program
func CleanUp(s ShaderComponent) {
	gl.DeleteProgram(s.id)
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
